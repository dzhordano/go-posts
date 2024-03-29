package service

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
)

type PostsService struct {
	repo repository.Posts

	commentsService Comments
}

func NewPostsService(repo repository.Posts, commentsService Comments) *PostsService {
	return &PostsService{
		repo:            repo,
		commentsService: commentsService,
	}
}

func (s *PostsService) Create(ctx context.Context, input domain.Post, userId uint) error {
	return s.repo.Create(ctx, input, userId)
}

func (s *PostsService) GetAll(ctx context.Context) ([]domain.Post, error) {
	return s.repo.GetAll(ctx)
}

func (s *PostsService) GetById(ctx context.Context, postId uint) (domain.Post, error) {
	return s.repo.GetById(ctx, postId)
}

func (s *PostsService) Update(ctx context.Context, input domain.UpdatePostInput, postId uint) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, input, postId)
}

func (s *PostsService) Delete(ctx context.Context, postId uint) error {
	return s.repo.Delete(ctx, postId)
}

func (s *PostsService) GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error) {
	return s.repo.GetAllUser(ctx, userId)
}

func (s *PostsService) DeleteUser(ctx context.Context, postId, userId uint) error {
	return s.repo.DeleteUser(ctx, postId, userId)
}

func (s *PostsService) UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateUser(ctx, input, postId, userId)
}

func (s *PostsService) Like(ctx context.Context, postId uint) error {
	return s.repo.Like(ctx, postId)
}

func (s *PostsService) RemoveLike(ctx context.Context, postId uint) error {
	return s.repo.RemoveLike(ctx, postId)
}

func (s *PostsService) Report(ctx context.Context, postId, userId uint) error {
	return s.repo.Report(ctx, postId, userId)
}

func (s *PostsService) GetAllReports(ctx context.Context) ([]domain.Report, error) {
	return s.repo.GetAllReports(ctx)
}
