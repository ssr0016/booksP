package migrator

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

type Migrator struct {
	engine       *xorm.Engine
	log          *zap.Logger
	Dialect      Dialect
	migrationIds map[string]struct{}
	migrations   []Migration
}

type MigrationLog struct {
	Id          int64
	MigrationID string `xorm:"migration_id"`
	SQL         string `xorm:"sql"`
	Success     bool
	Error       string
	Timestamp   time.Time
}

func NewMigrator(engine *xorm.Engine) *Migrator {
	mg := &Migrator{}
	mg.engine = engine
	mg.log = zap.L().Named("migrator")
	mg.migrations = make([]Migration, 0)
	mg.Dialect = NewDialect(mg.engine)
	mg.migrationIds = make(map[string]struct{})
	return mg
}

func (mg *Migrator) MigrationsCount() int {
	return len(mg.migrations)
}

func (mg *Migrator) AddMigration(id string, m Migration) {
	if _, ok := mg.migrationIds[id]; ok {
		panic(fmt.Sprintf("migration id conflict: %s", id))
	}

	m.SetId(id)
	mg.migrations = append(mg.migrations, m)
	mg.migrationIds[id] = struct{}{}
}

func (mg *Migrator) GetMigrationLog() (map[string]MigrationLog, error) {
	logMap := make(map[string]MigrationLog)
	logItems := make([]MigrationLog, 0)

	exists, err := mg.engine.IsTableExist(new(MigrationLog))
	if err != nil {
		return nil, fmt.Errorf("%v: %w", "failed to check table existence", err)
	}

	if !exists {
		return logMap, nil
	}

	if err = mg.engine.Find(&logItems); err != nil {
		return nil, err
	}

	for _, logItem := range logItems {
		if !logItem.Success {
			continue
		}
		logMap[logItem.MigrationID] = logItem
	}

	return logMap, nil
}

func (mg *Migrator) Start() error {
	mg.log.Info("starting DB migrations")

	logMap, err := mg.GetMigrationLog()
	if err != nil {
		return err
	}

	migrationsPerformed := 0
	migrationsSkipped := 0
	start := time.Now()
	for _, m := range mg.migrations {
		m := m
		_, exists := logMap[m.Id()]
		if exists {
			mg.log.Debug("skipping migration: Already executed",
				zap.String("id", m.Id()),
			)
			migrationsSkipped++
			continue
		}

		sql := m.SQL(mg.Dialect)

		record := MigrationLog{
			MigrationID: m.Id(),
			SQL:         sql,
			Timestamp:   time.Now(),
		}

		err := mg.inTransaction(func(sess *xorm.Session) error {
			err := mg.exec(m, sess)
			if err != nil {
				mg.log.Error("executing migration condition failed",
					zap.String("sql", sql),
					zap.Error(err),
				)

				record.Error = err.Error()
				if _, err := sess.Insert(&record); err != nil {
					return err
				}
				return err
			}
			record.Success = true
			_, err = sess.Insert(&record)
			if err == nil {
				migrationsPerformed++
			}
			return err
		})
		if err != nil {
			return fmt.Errorf("%v: %w", "migration failed", err)
		}
	}

	mg.log.Info("migrations completed",
		zap.Int("performed", migrationsPerformed),
		zap.Int("skipped", migrationsSkipped),
		zap.Duration("duration", time.Since(start)),
	)

	return mg.engine.Sync2()
}

func (mg *Migrator) exec(m Migration, sess *xorm.Session) error {

	mg.log.Info("executing migration",
		zap.String("id", m.Id()),
	)

	condition := m.GetCondition()
	if condition != nil {
		sql, args := condition.Sql(mg.Dialect)

		if sql != "" {
			mg.log.Debug("executing migration condition sql",
				zap.String("id", m.Id()),
				zap.String("sql", sql),
				// zap.ObjectValues("args", args),
			)

			results, err := sess.SQL(sql, args...).Query()
			if err != nil {
				mg.log.Error("executing migration condition failed",
					zap.String("id", m.Id()),
					zap.String("error", err.Error()),
				)
				return err
			}

			if !condition.IsFulfilled(results) {
				mg.log.Warn("skipping migration: Already executed, but not recorded in migration log",
					zap.String("id", m.Id()),
				)
				return nil
			}
		}
	}

	var err error
	if codeMigration, ok := m.(CodeMigration); ok {
		mg.log.Debug("Executing code migration",
			zap.String("id", m.Id()))

		err = codeMigration.Exec(sess, mg)
	} else {
		sql := m.SQL(mg.Dialect)
		mg.log.Debug("Executing sql migration",
			zap.String("id", m.Id()),
			zap.String("sql", sql))
		_, err = sess.Exec(sql)
	}

	if err != nil {
		mg.log.Error("Executing migration condition failed",
			zap.String("id", m.Id()),
			zap.Error(err),
		)
		return err
	}

	return nil
}

type dbTransactionFunc func(sess *xorm.Session) error

func (mg *Migrator) inTransaction(callback dbTransactionFunc) error {
	sess := mg.engine.NewSession()
	defer sess.Close()

	if err := sess.Begin(); err != nil {
		return err
	}

	if err := callback(sess); err != nil {
		if rollErr := sess.Rollback(); rollErr != nil {
			return fmt.Errorf("failed to roll back transaction due to error: %s", rollErr)
		}

		return err
	}

	if err := sess.Commit(); err != nil {
		return err
	}

	return nil
}
