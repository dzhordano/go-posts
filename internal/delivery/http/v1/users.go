package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
	}
}

func (h *Handler) userSignUp(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (h *Handler) userSignIn(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
