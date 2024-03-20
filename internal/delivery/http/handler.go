package http

import (
	"fmt"
	"net/http"

	"github.com/dzhordano/go-posts/docs"
	"github.com/dzhordano/go-posts/internal/common/config"
	v1 "github.com/dzhordano/go-posts/internal/delivery/http/v1"
	"github.com/dzhordano/go-posts/internal/service"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager ``
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
	// limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
	// corsMiddleware,
	)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Env != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}

	if cfg.Env != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "helow")
	})

	h.initApi(router)

	return router
}

func (h *Handler) initApi(router *gin.Engine) {
	v1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		v1.Init(api)
	}
}
