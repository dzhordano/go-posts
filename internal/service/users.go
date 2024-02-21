package service

import (
	"context"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) SignUP(ctx context.Context, input domain.UserSignUpInput) error {
	// TODO: implement password hashing here
	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Verification: domain.Verification{
			Code:       "no-code",
			IsVerified: false,
		},
		RegisteredAt: time.Now(),
		LastOnline:   time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIN(ctx context.Context, input domain.UserSignInInput) error {
	// TODO: implement
	panic("todo")
}
