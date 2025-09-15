package app

import (
	"context"
	"e-comm/config"
	"e-comm/internal/delivery"
	"e-comm/internal/repository"
	"e-comm/internal/server"
	"e-comm/internal/service"
	postgres "e-comm/pkg/adapter"
	"fmt"
)

func StartWebServer(ctx context.Context, cfg *config.Config) error {
	pg, err := postgres.New(ctx, &cfg.PG)
	if err != nil {
		return fmt.Errorf("could not init postgres connection: %w", err)
	}

	repos := repository.NewRepositories(pg)
	services := service.NewServices(service.Deps{
		Repos:  repos,
		Config: cfg,
	})

	handlers := delivery.NewHandler(services)

	srv := server.NewServer(cfg, handlers.Init(cfg))

	srv.Run()
	defer srv.Stop()
	fmt.Println("Web server started!")

	<-ctx.Done()

	return nil
}
