package http

import (
	"net/http"

	v1 "github.com/dzhordano/go-posts/internal/delivery/http/v1"
	"github.com/dzhordano/go-posts/internal/service"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

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
