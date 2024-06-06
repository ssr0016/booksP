package postgres

import (
	"practice/infra/storage/db"
	"practice/infra/storage/db/dbimpl"
	"practice/migrator"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

type DB interface {
	db.DB
}

type Tx = db.Tx

type Migrator interface {
	AddMigration(mg *migrator.Migrator)
}

type postgresdb struct {
	log        *zap.Logger
	engine     *xorm.Engine
	migrations Migrator
	dialect    migrator.Dialect
	db.DB
}

func New(migrations Migrator, connection string) (DB, error) {
	var err error

	p := &postgresdb{
		log: zap.L().Named("postgres"),
	}

	engine, err := xorm.NewEngine("postgres", connection)
	if err != nil {
		return nil, err
	}

	p.engine = engine
	p.migrations = migrations
	p.dialect = migrator.NewDialect(p.engine)

	engine.SetTZDatabase(time.UTC)

	if err = p.Migrate(); err != nil {
		p.log.Error("migration failed err: %v", zap.Any("errors", err))
		return nil, err
	}

	p.DB = dbimpl.NewSqlx(sqlx.NewDb(p.engine.DB().DB, p.GetDialect().DriverName()))
	return p, nil
}

func (p *postgresdb) Migrate() error {
	migrator := migrator.NewMigrator(p.engine)
	p.migrations.AddMigration(migrator)
	return migrator.Start()
}

func (p *postgresdb) GetDialect() migrator.Dialect {
	return p.dialect
}
