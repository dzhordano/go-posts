package main

import (
	"fmt"
	"log"

	"github.com/dzhordano/go-posts/internal/common/config"
	delivery "github.com/dzhordano/go-posts/internal/delivery/http"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/internal/service"
	"github.com/dzhordano/go-posts/pkg/postgres"
	"github.com/dzhordano/go-posts/pkg/server"
	_ "github.com/lib/pq"
)

var cfgPath = "configs"

func main() {

	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatal("error initializing config")
	}

	fmt.Printf("%s %s %s %s %s %s", cfg.PG.Host, cfg.PG.Port, cfg.PG.DBName, cfg.PG.Username, cfg.PG.Password, cfg.PG.SSLMode)

	db, err := postgres.NewClient(cfg.PG.Host, cfg.PG.Port, cfg.PG.DBName, cfg.PG.Username, cfg.PG.Password, cfg.PG.SSLMode)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// Init repositories
	repos := repository.NewRepos(db)

	// Init services
	services := service.NewService(service.Deps{
		Repos: repos,
	})

	// Init handlers
	handlers := delivery.NewHandler(services)

	// Init server
	srv := server.NewServer(handlers.Init())

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatal("failed to start server")
	}
}
