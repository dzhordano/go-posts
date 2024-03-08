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

	postsService Posts

	accessTokenTLL  time.Duration
	refreshTokenTLL time.Duration
}

func NewUsersService(repo repository.Users, hasher hasher.PassworsHasher, tokenManager auth.TokenManager, postsService Posts, attl, rttl time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		postsService:    postsService,
		tokenManager:    tokenManager,
		accessTokenTLL:  attl,
		refreshTokenTLL: rttl,
	}
}

func (s *UsersService) SignUP(ctx context.Context, input domain.UserSignUpInput) error {
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
		RegisteredAt: time.Now(),
		LastOnline:   time.Now(),
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIN(ctx context.Context, input domain.UserSignInInput) (Tokens, error) {
	// TODO: and also return tokens, not id
	passwordHash, err := s.hasher.GeneratePasswordHash(input.Password)
	if err != nil {
		return Tokens{}, err
	}
	input.Password = passwordHash

	user, err := s.repo.GetByCredentials(ctx, input)
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
