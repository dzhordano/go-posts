package main

import (
	"context"
	"log"

	"github.com/dzhordano/go-posts/internal/common/config"
	delivery "github.com/dzhordano/go-posts/internal/delivery/http"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/internal/service"
	"github.com/dzhordano/go-posts/pkg/postgres"
	"github.com/dzhordano/go-posts/pkg/server"
)

var cfgPath = "configs"

func main() {
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatal("error initializing config")
	}

	pgclient, err := postgres.NewClient(context.Background(), postgres.DBConfig{
		Username: cfg.PG.Username,
		Password: cfg.PG.Password,
		Host:     cfg.PG.Host,
		Port:     cfg.PG.Port,
		Database: cfg.PG.DBName,
		SSLMode:  cfg.PG.SSLMode,
		MaxAtts:  5,
	})

	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// Init repositories
	repos := repository.NewRepos(pgclient)

	// Init services
	services := service.NewService(service.Deps{
		Repos: repos,
	})

	// Init handlers
	handlers := delivery.NewHandler(services)

	// Init server
	srv := server.NewServer(cfg, handlers.Init())

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatal("failed to start server")
	}
}
