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

// configs dir path
var cfgPath = "configs"

// TODO: зач в структуре нужны поля типо сессии если я их не использую да и не сохраняю туда данные (чет я тупеньки)
// TODAY: простой фронт с HTMX
// NEXT: сделать эндпоинты для админов, чтобы админ могу манипулировать пользователями и постами и подумать над тем как "прикрепить" правильно
// комменты к посту
// NEXT: следующим надо внедрить комментарии к посту и дать проработать их удаление или изменение, также сделать suspend
// чтобы такие комменты заменялись на "suspended"
// Redis: после всего надо подумать что кэшировать (посты, мб токены, или ???)

func main() {
	// load cfg from env and yaml
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatal("error initializing config")
	}

	// init token manager (jwt tokens manager)
	tokenManager, err := auth.NewManager(cfg.AUTH.JWT.SigningKey)
	if err != nil {
		log.Fatalf("failed on: %v", err)
	}

	// hasher for hashing
	hasher := hasher.NewSHA256Hasher(cfg.AUTH.PasswordSalt)

	// launch postgres client and connect to pool
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
	// FIXME: dont forget to put pgxpool.Close() after graceful shutdown or somewhere where it belongs
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
