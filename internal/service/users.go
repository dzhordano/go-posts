package service

import (
	"context"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/google/uuid"
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
	// TODO: implement password hashing for input.password
	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Verification: domain.Verification{
			Code:     "no-code",
			Verified: false,
		},
		RegisteredAt: time.Now(),
		LastOnline:   time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIN(ctx context.Context, input domain.UserSignInInput) (uuid.UUID, error) {
	// TODO: implement password hashing for input.passowrd
	// 		 and also return tokens, not id
	user, err := s.repo.GetByCredentials(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return user.UID, nil
}
