package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"practice/infra/storage/postgres"
	"practice/library/books"
	"practice/library/config"
	"practice/routing"
	"time"

	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

type Server struct {
	Router          *routing.Router
	Cfg             *config.Config
	Dependencies    *Dependencies
	log             *zap.Logger
	ShutdownTimeout time.Duration
}

type Dependencies struct {
	Cfg        *config.Config
	PostgresDB postgres.DB
	BooksSvc   books.Service
}

var decoder = schema.NewDecoder()

func NewServer(deps *Dependencies, cfg *config.Config) *Server {
	r := routing.NewRouter()
	decoder.IgnoreUnknownKeys(true)

	return &Server{
		Dependencies:    deps,
		Router:          r,
		Cfg:             cfg,
		log:             zap.L().Named("rest.server"),
		ShutdownTimeout: time.Duration(60) * time.Second,
	}
}

func (s *Server) registerRoutes() {
	r := s.Router

	s.NewBooksHandler(r)
}

// func (s *Server) addMiddlewares() {
// 	r := s.Router
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.RequestTracing(s.Dependencies.Tracing))
// 	r.Use(middleware.Recovery)
// 	r.Use(s.Dependencies.ContextHandler.Middleware)
// }

func (s *Server) RunServer(server *http.Server, listener net.Listener, shutDownTimeout time.Duration, stopCh <-chan struct{}) (<-chan struct{}, <-chan struct{}, error) {
	// Shutdown server gracefully.
	serverShutdownCh, listenerStoppedCh := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(serverShutdownCh)
		<-stopCh
		ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
		server.Shutdown(ctx)
		cancel()
	}()

	go func() {
		defer close(listenerStoppedCh)

		if server.TLSConfig != nil {
			listener = tls.NewListener(listener, server.TLSConfig)
		}

		err := server.Serve(listener)

		msg := fmt.Sprintf("Stopped listening on %s", listener.Addr().String())
		select {
		case <-stopCh:
			s.log.Info(msg)
		default:
			panic(fmt.Sprintf("%s due to error: %v", msg, err))
		}
	}()

	return serverShutdownCh, listenerStoppedCh, nil
}

func (s *Server) Run(ctx context.Context) error {
	stopCh := ctx.Done()

	// s.addMiddlewares()
	s.registerRoutes()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Dependencies.Cfg.Server.HTTPPort))
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler: s.Router,
	}

	go func() {
		<-stopCh
		server.Shutdown(context.Background())
		s.log.Info("shutting down HTTP server")
	}()

	s.log.Info("starting HTTP server ", zap.String("port", s.Dependencies.Cfg.Server.HTTPPort))

	var stoppedCh <-chan struct{}
	var listenerStoppedCh <-chan struct{}
	stoppedCh, listenerStoppedCh, err = s.RunServer(server, listener, s.ShutdownTimeout, stopCh)
	if err != nil {
		return err
	}

	s.log.Info("[graceful-termination] waiting for shutdown to be initiated")
	<-stopCh

	<-listenerStoppedCh
	<-stoppedCh
	s.log.Info("[graceful-termination] server is exiting")

	return nil
}
