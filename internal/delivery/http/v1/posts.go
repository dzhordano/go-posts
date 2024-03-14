package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initPostsRoutes(api *gin.RouterGroup) {
	posts := api.Group("/posts")
	{
		posts.GET("/", h.getAllPosts)
		posts.GET("/:id", h.getPostById)

		posts.GET("/:id/comments", h.getPostComments)
	}
}

func (h *Handler) getAllPosts(c *gin.Context) {
	// TODO: think of retrieving random posts instead, so request is not long and make retrieval random
	posts, err := h.services.Posts.GetAll(c.Request.Context())
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: posts,
	})
}

func (h *Handler) getPostById(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.services.Posts.GetById(c.Request.Context(), uint(postId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: post,
	})
}

func (h *Handler) getPostComments(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	comments, err := h.services.Comments.GetComments(c.Request.Context(), uint(postId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: comments,
	})
}
