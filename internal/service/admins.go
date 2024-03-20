package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/hasher"
)

type AdminsService struct {
	repo         repository.Admins
	hasher       hasher.PassworsHasher
	tokenManager auth.TokenManager

	postsService Posts
	usersService Users

	accessTokenTLL  time.Duration
	refreshTokenTLL time.Duration
}

func NewAdminsService(repo repository.Admins, hasher hasher.PassworsHasher, tokenManager auth.TokenManager, postsService Posts, usersService Users, attl, rttl time.Duration) *AdminsService {
	return &AdminsService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		postsService:    postsService,
		usersService:    usersService,
		accessTokenTLL:  attl,
		refreshTokenTLL: rttl,
	}
}

func (s *AdminsService) SignIN(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.GeneratePasswordHash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	input.Password = passwordHash

	admin, err := s.repo.GetByCredentials(ctx, input.Email, input.Password)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, admin.ID)
}

func (s *AdminsService) createSession(ctx context.Context, adminId uint) (res Tokens, err error) {

	res.AccessToken, err = s.tokenManager.CreateJWT(fmt.Sprintf("%x", adminId), s.accessTokenTLL)
	if err != nil {
		return
	}

	res.RefreshToken, err = s.tokenManager.CreateJWT(fmt.Sprintf("%x", adminId), s.refreshTokenTLL)
	if err != nil {
		return
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTLL),
	}

	err = s.repo.CreateSession(ctx, adminId, session)

	return
}

func (s *AdminsService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	admin, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, admin.ID)
}

func (s *AdminsService) UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateUser(ctx, input, userId)
}

func (s *AdminsService) DeleteUser(ctx context.Context, userId uint) error {
	return s.repo.DeleteUser(ctx, userId)
}

func (s *AdminsService) SuspendUser(ctx context.Context, userId uint) error {
	return s.repo.SuspendUser(ctx, userId)
}

func (s *AdminsService) SuspendPost(ctx context.Context, postId uint) error {
	return s.repo.SuspendPost(ctx, postId)
}

func (s *AdminsService) CensorComment(ctx context.Context, commId uint) error {
	return s.repo.CensorComment(ctx, commId)
}

func (s *AdminsService) DeleteComment(ctx context.Context, commId uint) error {
	return s.repo.DeleteComment(ctx, commId)
}
