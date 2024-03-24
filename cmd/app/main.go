package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/dzhordano/go-posts/internal/common/config"
	delivery "github.com/dzhordano/go-posts/internal/delivery/http"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/internal/service"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/email/smtp"
	"github.com/dzhordano/go-posts/pkg/hasher"
	"github.com/dzhordano/go-posts/pkg/logger"
	"github.com/dzhordano/go-posts/pkg/otp"
	"github.com/dzhordano/go-posts/pkg/postgres"
	"github.com/dzhordano/go-posts/pkg/server"
)

// configs dir path
var cfgPath = "configs"

// TODO: do more user tests. admin+posts+comments tests.

//	@title			Go-Posts Api
//	@version		1.0
//	@description	...
//	@host			localhost:8081
//	@BasePath		/api/v1

//	@securityDefinitions.apiKey	UserAuth
//	@in							header
//	@name						Authorization

// @securityDefinitions.apiKey	AdminAuth
// @in							header
// @name						Authorization
func main() {
	// TODO: user this logger
	slog.SetDefault(logger.InitLogger())

	// load cfg from env and yaml
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatal("error initializing config")
	}

	// one time password* generator init
	otpGenerator := otp.NewGOTPGenerator()

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
	defer pgclient.Close()

	fmt.Println(cfg.AUTH.VerificationCodeLength)
	// Init email sender
	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Init repositories
	repos := repository.NewRepos(pgclient)

	// Init services
	services := service.NewService(service.Deps{
		Repos:        repos,
		Hasher:       hasher,
		TokenManager: tokenManager,
		EmailSender:  emailSender,
		OtpGenerator: otpGenerator,

		AccessTokenTTL:         cfg.AUTH.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.AUTH.JWT.RefreshTokenTTL,
		VerificationCodeLength: cfg.AUTH.VerificationCodeLength,
	})

	// Init handlers
	handlers := delivery.NewHandler(services, tokenManager)

	// Init server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatal("failed to start server")
	}
}
