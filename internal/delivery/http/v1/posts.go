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

//	@Summary		Get All Posts
//	@Tags			posts
//	@Description	get all  posts
//	@ID				get-all-posts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dataResponse
//	@Failure		404		{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/posts [get]
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

//	@Summary		Get Post By Id
//	@Tags			posts
//	@Description	get post by id
//	@ID				get-post-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"post id"
//	@Success		200		{object}	dataResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/posts/{id} [get]
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

//	@Summary		Get Post Comments
//	@Tags			posts
//	@Description	get post comments
//	@ID				get-post-comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"post id"
//	@Success		200		{object}	dataResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/posts/{id}/comments [get]
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
