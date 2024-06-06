package app

import (
	"context"
	"errors"
	"fmt"
	"practice/infra/storage/postgres"
	"practice/library/books/booksimpl"
	"practice/library/config"
	"practice/library/protocol/rest"
	"practice/library/storage/postgres/migrations"
	"practice/registry"
	"sync"

	"go.uber.org/zap"
)

type Server struct {
	log        *zap.Logger
	services   []registry.RunFunc
	postgresdb postgres.DB
}

func NewServer(isStandaloneMode bool) (*Server, error) {

	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	postgresdb, err := postgres.New(migrations.New(), cfg.Postgres.ConnectionString())
	if err != nil {
		return nil, err
	}

	// postgres store
	// bookStore := booksimpl.NewStore(postgresdb)

	bookSvc := booksimpl.NewService(postgresdb, cfg)

	restServer := rest.NewServer(&rest.Dependencies{
		Cfg:        cfg,
		PostgresDB: postgresdb,
		BooksSvc:   bookSvc,
	}, cfg)

	services := registry.NewServiceRegistry(
		restServer.Run,
	)

	if isStandaloneMode {
		services = registry.NewServiceRegistry(
			restServer.Run,
		)
	}

	return &Server{
		services:   services.GetServices(),
		postgresdb: postgresdb,
		log:        zap.L().Named("apiserver"),
	}, nil
}

func (s *Server) Run(ctx context.Context) {
	defer func() {
		s.postgresdb.Close()
	}()

	var wg sync.WaitGroup
	wg.Add(len(s.services))

	for _, svc := range s.services {
		go func(svc registry.RunFunc) error {
			defer wg.Done()
			err := svc(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {
				s.log.Error("stopped server", zap.String("service", serviceName), zap.Error(err))
				return fmt.Errorf("%s run error: %w", serviceName, err)
			}
			return nil
		}(svc)
	}

	wg.Wait()
}
