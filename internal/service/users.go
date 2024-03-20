package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/hasher"
)

type UsersService struct {
	repo         repository.Users
	hasher       hasher.PassworsHasher
	tokenManager auth.TokenManager

	postsService    Posts
	commentsService Comments

	accessTokenTLL  time.Duration
	refreshTokenTLL time.Duration
}

func NewUsersService(repo repository.Users, hasher hasher.PassworsHasher, tokenManager auth.TokenManager, postsService Posts, commentsService Comments, attl, rttl time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		postsService:    postsService,
		commentsService: commentsService,
		tokenManager:    tokenManager,
		accessTokenTLL:  attl,
		refreshTokenTLL: rttl,
	}
}

func (s *UsersService) SignUP(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.GeneratePasswordHash(input.Password)
	if err != nil {
		return err
	}

	// TODO: change verification data when implementing verification
	user := domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: passwordHash,
		Verification: domain.Verification{
			Code:     "no-code",
			Verified: true,
		},
		Session: domain.Session{
			RefreshToken: "null",
			ExpiresAt:    time.Now(),
		},
		RegisteredAt: time.Now(),
		LastOnline:   time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIN(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.GeneratePasswordHash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	input.Password = passwordHash

	user, err := s.repo.GetByCredentials(ctx, input.Email, input.Password)
	if err != nil {
		return Tokens{}, err
	}

	if !user.Verification.Verified {
		return Tokens{}, errors.New("user is not verified")
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) createSession(ctx context.Context, userId uint) (res Tokens, err error) {

	res.AccessToken, err = s.tokenManager.CreateJWT(fmt.Sprintf("%x", userId), s.accessTokenTLL)
	if err != nil {
		return
	}

	res.RefreshToken, err = s.tokenManager.CreateJWT(fmt.Sprintf("%x", userId), s.refreshTokenTLL)
	if err != nil {
		return
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTLL),
	}

	err = s.repo.CreateSession(ctx, userId, session)

	return
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) GetAll(ctx context.Context) ([]domain.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UsersService) GetById(ctx context.Context, userId uint) (domain.User, error) {
	return s.repo.GetById(ctx, userId)
}
