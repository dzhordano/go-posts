package v1

import (
	"net/http"
	"strconv"

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
			posts := auth.Group("/posts")
			{
				posts.GET("/", h.getUserPosts)
				posts.GET("/:id", h.getUserPostById) // FIXME: DELETE (USELESS) ??
				posts.GET("/:id/comments", h.getUserPostComments)

				posts.POST("/", h.createUserPost)
				posts.POST("/:id/comments", h.createPostComment)

				posts.PUT("/:id", h.updateUserPost)

				posts.DELETE("/:id", h.deleteUserPost)
				// posts.DELETE("/:id/comments/:cid") // TODO: do so user can delete his own comments
			}

			comments := auth.Group("/comments")
			{
				comments.GET("/", h.getUserComments)
			}
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

func (h *Handler) getUserPosts(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	posts, err := h.services.Posts.GetAllUser(c.Request.Context(), userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: posts,
	})
}

func (h *Handler) getUserPostById(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	post, err := h.services.Posts.GetByIdUser(c.Request.Context(), uint(postId), userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: post,
	})
}

func (h *Handler) createUserPost(c *gin.Context) {
	var input domain.Post

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	user, err := h.services.Users.GetById(c.Request.Context(), userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	input.Author = user.Name

	if err := h.services.Posts.Create(c.Request.Context(), input, userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "created",
	})
}

func (h *Handler) updateUserPost(c *gin.Context) {
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

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Posts.UpdateUser(c.Request.Context(), input, uint(postId), userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "updated",
	})
}

func (h *Handler) deleteUserPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Posts.DeleteUser(c.Request.Context(), uint(postId), userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "deleted",
	})
}

func (h *Handler) createPostComment(c *gin.Context) {
	var input domain.Comment
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	input.AuthorId = userId

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Comments.Create(c.Request.Context(), input, uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "created",
	})
}

func (h *Handler) getUserPostComments(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	comments, err := h.services.Comments.GetUserPostComments(c.Request.Context(), userId, uint(postId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: comments,
	})
}

func (h *Handler) getUserComments(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	comments, err := h.services.Comments.GetUserComments(c.Request.Context(), userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: comments,
	})
}
