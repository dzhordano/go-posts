package v1

import (
	"net/http"

	"github.com/dzhordano/go-posts/internal/domain"
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
	var input domain.UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.services.Users.SignUP(c, input); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) userSignIn(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
