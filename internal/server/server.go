package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"e-comm/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Http.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error occurred while running http server", slog.Any("error", err))
		}
	}()
}

func (s *Server) Stop() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Stopping server failed", slog.Any("error", err))
	}
}
