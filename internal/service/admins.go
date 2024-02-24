package service

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/google/uuid"
)

type AdminsService struct {
	repo repository.Admins
}

func NewAdminsService(repo repository.Admins) *AdminsService {
	return &AdminsService{
		repo: repo,
	}
}

func (s *AdminsService) SignIN(ctx context.Context, input domain.UserSignInInput) (uuid.UUID, error) {
	// TODO: implement hashing for input.password
	admin, err := s.repo.GetByCredentials(ctx, input)
	if err != nil {
		return uuid.Nil, err
	}

	return admin.UID, nil
}
