package main

import (
	"context"
	"e-comm/config"
	"e-comm/internal/app"
	"os"
	"os/signal"
	"syscall"

	"log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %w", err)
	}

	if err := app.StartWebServer(ctx, cfg); err != nil {
		log.Fatalf("failed to start server: %w", err)
	}
}
