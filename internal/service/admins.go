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

	accessTokenTLL  time.Duration
	refreshTokenTLL time.Duration
}

func NewAdminsService(repo repository.Admins, hasher hasher.PassworsHasher, tokenManager auth.TokenManager, postsService Posts, attl, rttl time.Duration) *AdminsService {
	return &AdminsService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		postsService:    postsService,
		accessTokenTLL:  attl,
		refreshTokenTLL: rttl,
	}
}

func (s *AdminsService) SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.GeneratePasswordHash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	input.Password = passwordHash
	fmt.Println(passwordHash)

	admin, err := s.repo.GetByCredentials(ctx, input)
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
