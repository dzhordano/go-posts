package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// empeti
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "helow")
	})

	return router
}
