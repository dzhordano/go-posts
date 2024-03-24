package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/service"
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
			auth.POST("/verify/:code", h.userVerify)

			posts := auth.Group("/posts")
			{
				posts.GET("/", h.getUserPosts)
				posts.POST("/", h.createUserPost)
				posts.PUT("/:id", h.updateUserPost)
				posts.DELETE("/:id", h.deleteUserPost)

				posts.GET("/:id/comments", h.getUserPostComments)
				posts.POST("/:id/comment", h.createPostComment)

				posts.POST("/:id/like", h.userLikePost)
				posts.POST("/:id/unlike", h.userRemoveLike)
				// posts.POST("/:id/react", h.userReactPost) TODO: later...
			}

			comments := auth.Group("/comments")
			{
				comments.GET("/", h.getUserComments)
				comments.PUT("/:id", h.updateUserComment)
				comments.DELETE("/:id", h.deleteUserComment)
			}
		}
	}
}

type userSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type userSignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

//	@Summary		User Verify
//	@Security		UserAuth
//	@Tags			users
//	@Description	verification endpoint for users
//	@ID				user-verify
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"verification code"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/users/verify/{code} [post]
func (h *Handler) userVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	id, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Users.Verify(c.Request.Context(), id, code); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{"success"})
}

//	@Summary		Sign Up
//	@Tags			users
//	@Description	registration endpoint for users
//	@ID				user-signup
//	@Accept			json
//	@Produce		json
//	@Param			input	body		userSignUpInput	true	"account info"
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/users/sign-up [post]
func (h *Handler) userSignUp(c *gin.Context) {
	var input userSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Users.SignUP(c.Request.Context(), service.UserSignUpInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "verification link sent",
	})
}

//	@Summary		Sign In
//	@Tags			users
//	@Description	login for users
//	@ID				user-signin
//	@Accept			json
//	@Produce		json
//	@Param			input	body		userSignInInput	true	"account credentials"
//	@Success		200		{object}	tokenResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/users/sign-in [post]
func (h *Handler) userSignIn(c *gin.Context) {
	var input userSignInInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	res, err := h.services.Users.SignIN(c.Request.Context(), service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
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

//	@Summary		Refresh Tokens
//	@Tags			users
//	@Description	refresh user's tokens
//	@ID				user-refresh-tokens
//	@Accept			json
//	@Produce		json
//	@Param			input	body		refreshInput	true	"refresh token"
//	@Success		200		{object}	tokenResponse
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/users/auth/refresh [post]
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

//	@Summary		Get User Posts
//	@Security		UserAuth
//	@Tags			users
//	@Description	get all user's posts
//	@ID				user-get-posts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dataResponse
//	@Failure		404		{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/users/posts [get]
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

type createUserPostInput struct {
	Title       string `json:"title" binding:"required,min=1"`
	Description string `json:"description" binding:"required,min=1"`
}

//	@Summary		Create User Post
//	@Security		UserAuth
//	@Tags			users
//	@Description	create post by user
//	@ID				user-create-post
//	@Param			input	body	createUserPostInput	true	"create user"
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response
//	@Failure		400,404	{object}	response
//	@Failure		500		{object}	response
//	@Failure		default	{object}	response
//	@Router			/users/posts/{id} [post]
func (h *Handler) createUserPost(c *gin.Context) {
	var input createUserPostInput

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

	if err := h.services.Posts.Create(c.Request.Context(), domain.Post{
		Title:       input.Title,
		Description: input.Description,
		Author:      user.Name,
		Suspended:   false,
		Created:     time.Now(),
		Updated:     time.Now(),
	}, userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "created",
	})
}

type updatePostInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

//	@Summary		Update User Post
//	@Security		UserAuth
//	@Tags			users
//	@Description	update user post
//	@ID				user-update-post
//	@Accept			json
//	@Produce		json
//	@Param			id					path		string			true	"post id"
//	@Param			input				body		updatePostInput	true	"update user post"
//	@Success		200					{object}	response
//	@Failure		400,404				{object}	response
//	@Failure		500					{object}	response
//	@Failure		default				{object}	response
//	@Router			/users/posts/{id} 	[put]
func (h *Handler) updateUserPost(c *gin.Context) {
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

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Posts.UpdateUser(c.Request.Context(), domain.UpdatePostInput{
		Title:       input.Title,
		Description: input.Description,
	}, uint(postId), userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "updated",
	})
}

//	@Summary		Delete User Post
//	@Security		UserAuth
//	@Tags			users
//	@Description	delete user post
//	@ID				user-delete-post
//	@Param			id	path	string	true	"post id"
//	@Accept			json
//	@Produce		json
//	@Success		200					{object}	response
//	@Failure		400,404				{object}	response
//	@Failure		500					{object}	response
//	@Failure		default				{object}	response
//	@Router			/users/posts/{id} 	[delete]
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

type createPostCommentInput struct {
	Data string `json:"data" binding:"required,min=1"`
}

//	@Summary		Create Post Comment
//	@Security		UserAuth
//	@Tags			users
//	@Description	create post comment
//	@ID				user-post-comment
//	@Accept			json
//	@Produce		json
//	@Param			id							path		string					true	"post id"
//	@Param			input						body		createPostCommentInput	true	"create post comment"
//	@Success		200							{object}	response
//	@Failure		400,404						{object}	response
//	@Failure		500							{object}	response
//	@Failure		default						{object}	response
//	@Router			/users/posts/{id}/comment 	[post]
func (h *Handler) createPostComment(c *gin.Context) {
	var input createPostCommentInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

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

	if err := h.services.Comments.Create(c.Request.Context(), domain.Comment{
		PostId:   uint(postId),
		AuthorId: userId,
		Data:     input.Data,
		Created:  time.Now(),
		Updated:  time.Now(),
	}, uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "created",
	})
}

//	@Summary		Get User's Post Comments
//	@Security		UserAuth
//	@Tags			users
//	@Description	get all user's post comments
//	@ID				user-get-post-comments
//	@Accept			json
//	@Produce		json
//	@Param			id							path		string					true	"post id"
//	@Param			input						body		createPostCommentInput	true	"create post comment"
//	@Success		200							{object}	dataResponse
//	@Failure		400,404						{object}	response
//	@Failure		500							{object}	response
//	@Failure		default						{object}	response
//	@Router			/users/posts/{id}/comments 	[get]
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

//	@Summary		Get All User Comments
//	@Security		UserAuth
//	@Tags			users
//	@Description	get all user comments
//	@ID				user-get-all-comments
//	@Accept			json
//	@Produce		json
//	@Success		200					{object}	dataResponse
//	@Failure		404					{object}	response
//	@Failure		500					{object}	response
//	@Failure		default				{object}	response
//	@Router			/users/comments/ 	[get]
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

//	@Summary		Update User Comment
//	@Security		UserAuth
//	@Tags			users
//	@Description	update user comment
//	@ID				update-user-comment
//	@Accept			json
//	@Produce		json
//	@Param			id						path		string	true	"comment id"
//	@Success		200						{object}	response
//	@Failure		400,404					{object}	response
//	@Failure		500						{object}	response
//	@Failure		default					{object}	response
//	@Router			/users/comments/{id} 	[put]
func (h *Handler) updateUserComment(c *gin.Context) {
	var input domain.UpdateCommentInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Comments.UpdateUser(c.Request.Context(), input, uint(commentId), userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "updated",
	})
}

//	@Summary		Delete User Comment
//	@Security		UserAuth
//	@Tags			users
//	@Description	delete user comment
//	@ID				delete-user-comment
//	@Accept			json
//	@Produce		json
//	@Param			id						path		string	true	"comment id"
//	@Success		200						{object}	response
//	@Failure		400,404					{object}	response
//	@Failure		500						{object}	response
//	@Failure		default					{object}	response
//	@Router			/users/comments/{id} 	[delete]
func (h *Handler) deleteUserComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	userId, err := h.getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Comments.DeleteUser(c.Request.Context(), uint(commentId), userId); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "deleted",
	})
}

//	@Summary		User Like Post
//	@Security		UserAuth
//	@Tags			users
//	@Description	user like post
//	@ID				user-like-post
//	@Accept			json
//	@Produce		json
//	@Param			id						path		string	true	"post id"
//	@Success		200						{object}	response
//	@Failure		400,404					{object}	response
//	@Failure		500						{object}	response
//	@Failure		default					{object}	response
//	@Router			/users/posts/{id}/like 	[post]
func (h *Handler) userLikePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Posts.Like(c.Request.Context(), uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "success",
	})
}

//	@Summary		User Unlike Post
//	@Security		UserAuth
//	@Tags			users
//	@Description	user unlike post
//	@ID				user-unlike-post
//	@Accept			json
//	@Produce		json
//	@Param			id							path		string	true	"post id"
//	@Success		200							{object}	response
//	@Failure		400,404						{object}	response
//	@Failure		500							{object}	response
//	@Failure		default						{object}	response
//	@Router			/users/posts/{id}/unlike 	[post]
func (h *Handler) userRemoveLike(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Posts.RemoveLike(c.Request.Context(), uint(postId)); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{
		Message: "success",
	})
}

// HERE TOO
// TODO: later...
// func (h *Handler) userReactPost(c *gin.Context) {
// 	panic("TODO")
// }
