package v1

import (
	"net/http"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdminsRoutes(api *gin.RouterGroup) {
	admins := api.Group("/admins")
	{
		admins.POST("/sign-in", h.adminSignIn)
	}
}

func (h *Handler) adminSignIn(c *gin.Context) {
	var input domain.UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Admins.SignIN(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, id)
}
