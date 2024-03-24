package v1

import (
	"net/http"
	"strconv"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/service"
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
				users.GET("/:id/posts", h.adminGetUserPosts)
				users.GET("/:id/comments", h.adminGetUserComments)

				users.POST("/:id/suspend", h.adminSuspendUser)

				users.PUT("/:id", h.adminAlterUser)

				users.DELETE("/:id", h.adminDeleteUser)
			}

			posts := auth.Group("/posts")
			{
				posts.POST("/:id/suspend", h.adminSuspendPost)

				posts.PUT("/:id", h.adminAlterPost)

				posts.DELETE("/:id", h.adminDeletePost)
			}

			comments := auth.Group("/comments")
			{
				comments.GET("/", h.adminGetComments)
				comments.POST("/", h.adminCreateComment)
				comments.PUT("/:id", h.adminUpdateComment)
				comments.DELETE("/:id", h.adminDeleteComment)
				comments.POST("/:id/censor", h.adminCensorComment)
			}
		}
	}
}

//	@Summary		Sign In
//	@Tags			admins
//	@Description	login for admin
//	@ID				admin-signup
//	@Accept			json
//	@Produce		json
//	@Param			input	body		userSignInInput	true	"account info"
//	@Success		200		{object}	tokenResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/sign-in [post]
func (h *Handler) adminSignIn(c *gin.Context) {
	var input userSignInInput

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tokens, err := h.services.Admins.SignIN(c.Request.Context(), service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

//	@Summary		Refresh Tokens
//	@Tags			admins
//	@Description	refresh admin's tokens
//	@ID				admin-refresh-tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body		refreshInput	true	"refresh token"
//	@Success		200		{object}	tokenResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admin/auth/refresh [post]
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

//	@Summary		Get Users
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	get all users
//	@ID				admin-get-users
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dataResponse
//	@Failure		404		{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users [get]
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

//	@Summary		Get User By Id
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	get user by id
//	@ID				admin-get-user-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"user id"
//	@Success		200		{object}	dataResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users/{id} [get]
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

//	@Summary		Alter User
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	alter user
//	@ID				admin-alter-user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"user id"
//	@Param			input	body		domain.UpdateUserInput	true	"update info"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users/{id} [put]
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

//	@Summary		Delete User
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	delete user
//	@ID				admin-delete-user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"user id"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users/{id} [delete]
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

//	@Summary		Suspend User
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	suspend user
//	@ID				admin-suspend-user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"user id"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users/{id}/suspend [post]
func (h *Handler) adminSuspendUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Admins.SuspendUser(c.Request.Context(), uint(userId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "suspended",
	})
}

//	@Summary		Alter Post
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	alter post
//	@ID				admin-alter-post
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string			true	"post id"
//	@Param			input	body		updatePostInput	true	"update info"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/posts/{id} [put]
func (h *Handler) adminAlterPost(c *gin.Context) {
	var input updatePostInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Posts.Update(c.Request.Context(), domain.UpdatePostInput{
		Title:       input.Title,
		Description: input.Description,
	}, uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "updated",
	})
}

//	@Summary		Delete Post
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	delete post
//	@ID				admin-delete-post
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"post id"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/posts/{id} [delete]
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

//	@Summary		Suspend Post
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	suspend post
//	@ID				admin-suspend-post
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"post id"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/posts/{id}/suspend [post]
func (h *Handler) adminSuspendPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Admins.SuspendPost(c.Request.Context(), uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response{
		Message: "suspended",
	})
}

//	@Summary		Censor Comment
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	censor comment
//	@ID				admin-censor-comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"comment id"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/comments/{id}/censor [post]
func (h *Handler) adminCensorComment(c *gin.Context) {
	commId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Admins.CensorComment(c.Request.Context(), uint(commId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "censored",
	})
}

//	@Summary		Delete Comment
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	delete comment
//	@ID				admin-delete-comment
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"comment id"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/comments/{id} [delete]
func (h *Handler) adminDeleteComment(c *gin.Context) {
	commId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Admins.DeleteComment(c.Request.Context(), uint(commId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "deleted",
	})
}

//	@Summary		Get User Posts
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	get all user's posts
//	@ID				admin-get-user-posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"user id"
//	@Success		200		{object}	dataResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users/{id}/posts [get]
func (h *Handler) adminGetUserPosts(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	posts, err := h.services.Posts.GetAllUser(c.Request.Context(), uint(userId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: posts,
	})
}

//	@Summary		Get User Comments
//	@Security		AdminAuth
//	@Tags			admins
//	@Description	get all user's comments
//	@ID				admin-get-user-comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"user id"
//	@Success		200		{object}	dataResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/admins/users/{id}/comments [get]
func (h *Handler) adminGetUserComments(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	comments, err := h.services.Comments.GetUserComments(c.Request.Context(), uint(userId))
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data: comments,
	})
}

// TODO: implement...
func (h *Handler) adminGetComments(c *gin.Context) {
	panic("TODO")
}

func (h *Handler) adminCreateComment(c *gin.Context) {
	panic("TODO")
}

func (h *Handler) adminUpdateComment(c *gin.Context) {
	panic("TODO")
}
