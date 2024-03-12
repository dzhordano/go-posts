package service

import (
	"context"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/hasher"
)

type (
	Users interface {
		SignUP(ctx context.Context, input domain.UserSignUpInput) error
		SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
		GetAll(ctx context.Context) ([]domain.User, error)
		GetById(ctx context.Context, userId uint) (domain.User, error)
	}

	Admins interface {
		SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)

		UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error
		DeleteUser(ctx context.Context, userId uint) error

		SuspendUser(ctx context.Context, userId uint) error
		SuspendPost(ctx context.Context, postId uint) error

		CensorComment(ctx context.Context, postId, commId uint) error
	}

	Posts interface {
		Create(ctx context.Context, input domain.Post, userId uint) error
		GetAll(ctx context.Context) ([]domain.Post, error)
		GetById(ctx context.Context, postId uint) (domain.Post, error)
		Update(ctx context.Context, input domain.UpdatePostInput, postId uint) error
		Delete(ctx context.Context, postId uint) error

		GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error)
		GetByIdUser(ctx context.Context, postId, userId uint) (domain.Post, error)
		UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error
		DeleteUser(ctx context.Context, postId, userId uint) error
	}

	Comments interface {
		Create(ctx context.Context, input domain.Comment, postId uint) error
		GetComments(ctx context.Context, postId uint) ([]domain.Comment, error)
		Delete(ctx context.Context, postId, commId uint) error

		GetUserComments(ctx context.Context, userId uint) ([]domain.Comment, error)
		GetUserPostComments(ctx context.Context, userId, postId uint) ([]domain.Comment, error)
	}
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Services struct {
	Users    Users
	Admins   Admins
	Posts    Posts
	Comments Comments
}

// services dependencies
type Deps struct {
	Repos        *repository.Repos
	Hasher       hasher.PassworsHasher
	TokenManager auth.TokenManager

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewService(deps Deps) *Services {
	commentsService := NewCommentsService(deps.Repos.Comments)
	postsService := NewPostsService(deps.Repos.Posts, commentsService)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, postsService, commentsService, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	adminsService := NewAdminsService(deps.Repos.Admins, deps.Hasher, deps.TokenManager, postsService, usersService, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users:    usersService,
		Admins:   adminsService,
		Posts:    postsService,
		Comments: commentsService,
	}
}
