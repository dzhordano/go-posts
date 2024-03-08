package service

import (
	"context"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/hasher"
)

type Users interface {
	SignUP(ctx context.Context, input domain.UserSignUpInput) error
	SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Admins interface {
	SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Posts interface {
	Create(ctx context.Context, input domain.Post, userId uint) error
	GetAll(ctx context.Context) ([]domain.Post, error)
	GetById(ctx context.Context, postId uint) (domain.Post, error)
	Update(ctx context.Context, input domain.UpdatePostInput) (domain.Post, error)
	Delete(ctx context.Context, postId uint) error
	// TODO: do i need to keep those here (i think yes)
	GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error)
	GetByIdUser(ctx context.Context, postId, userId uint) (domain.Post, error)
	UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error
	DeleteUser(ctx context.Context, postId, userId uint) error
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Services struct {
	Users  Users
	Admins Admins
	Posts  Posts
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
	postsService := NewPostsService(deps.Repos.Posts)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, postsService, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	adminsService := NewAdminsService(deps.Repos.Admins, deps.Hasher, deps.TokenManager, postsService, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users:  usersService,
		Admins: adminsService,
		Posts:  postsService,
	}
}
