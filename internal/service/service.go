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
	// TODO: change to tokens return after, not uuid
	SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Admins interface {
	// TODO: also to return tokens instead of uuid
	SignIN(ctx context.Context, input domain.UserSignInInput) (uint, error)
}

type Posts interface {
	Create(ctx context.Context, input domain.Post, userId uint) error
	GetAll(ctx context.Context) ([]domain.Post, error)
	GetById(ctx context.Context, postId uint) (domain.Post, error)
	Update(ctx context.Context, input domain.UpdatePostInput) (domain.Post, error)
	Delete(ctx context.Context) error
	// TODO: do i need to keep those here (i think yes)
	GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error)
	GetByIdUser(ctx context.Context, postId, userId uint) (domain.Post, error)
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
	adminsService := NewAdminsService(deps.Repos.Admins)

	return &Services{
		Users:  usersService,
		Admins: adminsService,
		Posts:  postsService,
	}
}
