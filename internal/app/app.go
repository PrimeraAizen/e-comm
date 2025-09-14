package app

import (
	"context"
	"e-comm/config"
	"e-comm/internal/delivery"
	"e-comm/internal/repository"
	"e-comm/internal/server"
	"e-comm/internal/service"
	"fmt"
)

func StartWebServer(ctx context.Context, cfg *config.Config) error {
	repos := repository.NewRepositories()
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
