package service

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
)

type Users interface {
	SignUP(ctx context.Context, input domain.UserSignUpInput) error
	SignIN(ctx context.Context, input domain.UserSignInInput) error
}

type Admins interface {
	SignIN(ctx context.Context, input domain.UserSignInInput) error
}

type Posts interface {
	Create(ctx context.Context, title, description string) error
	GetAll(ctx context.Context) ([]*domain.Post, error)
	GetById(ctx context.Context, postId uint) (*domain.Post, error)
	Update(ctx context.Context, input domain.UpdatePostInput) (*domain.Post, error)
	Delete(ctx context.Context) error
}

type Services struct {
	Users  Users
	Admins Admins
	Posts  Posts
}

func NewService() *Services {
	return &Services{}
}
