package v1

import (
	"net/http"
	"strconv"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdminsRoutes(api *gin.RouterGroup) {
	admins := api.Group("/admins")
	{
		admins.POST("/sign-in", h.adminSignIn)
		admins.POST("/auth/refresh", h.adminRefresh)

		auth := admins.Group("/", h.adminIdentity)
		{
			users := auth.Group("/users")
			{
				users.GET("/", h.adminGetUsers)
				users.GET("/:id", h.adminGetUserById)
				users.PUT("/:id", h.adminAlterUser)
				users.DELETE("/:id", h.adminDeleteUser)
			}

			posts := auth.Group("/posts")
			{
				posts.GET("/", h.adminGetPosts)
				posts.GET("/:id", h.adminGetPostById)
				posts.PUT("/:id", h.adminAlterPost)
				posts.DELETE("/:id", h.adminDeletePost)
			}
		}
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

func (h *Handler) adminRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	res, err := h.services.Admins.RefreshTokens(c.Request.Context(), inp.RToken)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) adminGetUsers(c *gin.Context) {
	users, err := h.services.Users.GetAll(c.Request.Context())
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: users,
	})
}

func (h *Handler) adminGetUserById(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	user, err := h.services.Users.GetById(c.Request.Context(), uint(userId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: user,
	})
}

func (h *Handler) adminAlterUser(c *gin.Context) {
	var input domain.UpdateUserInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Admins.UpdateUser(c.Request.Context(), input, uint(userId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "updated",
	})
}

func (h *Handler) adminDeleteUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Admins.DeleteUser(c.Request.Context(), uint(userId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "deleted",
	})
}

func (h *Handler) adminGetPosts(c *gin.Context) {
	posts, err := h.services.Posts.GetAll(c.Request.Context())
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: posts,
	})
}

func (h *Handler) adminGetPostById(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	posts, err := h.services.Posts.GetById(c.Request.Context(), uint(postId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: posts,
	})
}

func (h *Handler) adminAlterPost(c *gin.Context) {
	var input domain.UpdatePostInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Posts.Update(c.Request.Context(), input, uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "updated",
	})
}

func (h *Handler) adminDeletePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Posts.Delete(c.Request.Context(), uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response{
		Message: "deleted",
	})
}
