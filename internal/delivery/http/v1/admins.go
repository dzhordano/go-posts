package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdminsRoutes(api *gin.RouterGroup) {
	admins := api.Group("/admins")
	{
		admins.POST("/sign-in", h.adminSignIn)
	}
}

func (h *Handler) adminSignIn(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
