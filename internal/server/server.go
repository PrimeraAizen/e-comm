package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/PrimeraAizen/e-comm/config"
	"github.com/PrimeraAizen/e-comm/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	logger     *logger.Logger
}

func NewServer(cfg *config.Config, handler http.Handler, appLogger *logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              net.JoinHostPort(cfg.Http.Host, cfg.Http.Port),
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: appLogger,
	}
}

func (s *Server) Run() {
	go func() {
		s.logger.WithComponent("server").Info("HTTP server listening")
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.WithComponent("server").WithError(err).Error("HTTP server error")
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.WithComponent("server").Info("Initiating graceful shutdown")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.WithComponent("server").WithError(err).Error("Error during server shutdown")
		return err
	}

	s.logger.WithComponent("server").Info("HTTP server stopped gracefully")
	return nil
}
