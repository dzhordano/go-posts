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
		users.POST("/auth/refresh", h.userRefresh)

		auth := users.Group("/", h.userIdentity)
		{
			auth.GET("/test", func(c *gin.Context) { c.JSON(http.StatusOK, "Hello, world!") })
			// posts := auth.Group("/posts")
		}
	}
}

func (h *Handler) userSignUp(c *gin.Context) {
	var input domain.UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Users.SignUP(c.Request.Context(), input); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) userSignIn(c *gin.Context) {
	var input domain.UserSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.services.Users.SignIN(c.Request.Context(), input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		res.AccessToken,
		res.RefreshToken,
	})
}

type refreshInput struct {
	RToken string `json:"token" binding:"required"`
}

func (h *Handler) userRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Users.RefreshTokens(c.Request.Context(), inp.RToken)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
