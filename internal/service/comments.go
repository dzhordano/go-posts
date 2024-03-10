package service

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
)

type CommentsService struct {
	repo repository.Comments
}

func NewCommentsService(repo repository.Comments) *CommentsService {
	return &CommentsService{
		repo: repo,
	}
}

func (s *CommentsService) Create(ctx context.Context, input domain.Comment, postId uint) error {
	return s.repo.Create(ctx, input, postId)
}

func (s *CommentsService) GetComments(ctx context.Context, postId uint) ([]domain.Comment, error) {
	return s.repo.GetComments(ctx, postId)
}
