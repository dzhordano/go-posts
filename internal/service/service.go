package service

import (
	"context"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/hasher"
	"github.com/google/uuid"
)

type Users interface {
	SignUP(ctx context.Context, input domain.UserSignUpInput) error
	// TODO: change to tokens return after, not uuid
	SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Admins interface {
	// TODO: also to return tokens instead of uuid
	SignIN(ctx context.Context, input domain.UserSignInInput) (uuid.UUID, error)
}

type Posts interface {
	Create(ctx context.Context, title, description string) error
	GetAll(ctx context.Context) ([]domain.Post, error)
	GetById(ctx context.Context, postId uint) (domain.Post, error)
	Update(ctx context.Context, input domain.UpdatePostInput) (domain.Post, error)
	Delete(ctx context.Context) error
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
	return &Services{
		Users:  NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Admins: NewAdminsService(deps.Repos.Admins),
	}
}
