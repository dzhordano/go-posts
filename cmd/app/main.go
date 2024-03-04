package main

import (
	"context"
	"log"

	"github.com/dzhordano/go-posts/internal/common/config"
	delivery "github.com/dzhordano/go-posts/internal/delivery/http"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/internal/service"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/hasher"
	"github.com/dzhordano/go-posts/pkg/postgres"
	"github.com/dzhordano/go-posts/pkg/server"
)

var cfgPath = "configs"

// TODO: зач в структуре нужны поля типо сессии если я их не использую да и не сохраняю туда данные (чет я тупеньки)

func main() {
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatal("error initializing config")
	}

	tokenManager, err := auth.NewManager(cfg.AUTH.JWT.SigningKey)
	if err != nil {
		log.Fatalf("failed on: %v", err)
	}

	hasher := hasher.NewSHA256Hasher(cfg.AUTH.PasswordSalt)

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
	// TODO: dont forget to put pgxpool.Close() after graceful shutdown or somewhere where it belongs
	repos := repository.NewRepos(pgclient)

	// Init services
	services := service.NewService(service.Deps{
		Repos:        repos,
		Hasher:       hasher,
		TokenManager: tokenManager,

		AccessTokenTTL:  cfg.AUTH.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.AUTH.JWT.RefreshTokenTTL,
	})

	// Init handlers
	handlers := delivery.NewHandler(services, tokenManager)

	// Init server
	srv := server.NewServer(cfg, handlers.Init())

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatal("failed to start server")
	}
}
