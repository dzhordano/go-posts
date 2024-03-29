package service

import (
	"context"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/dzhordano/go-posts/internal/repository"
	"github.com/dzhordano/go-posts/pkg/auth"
	"github.com/dzhordano/go-posts/pkg/email"
	"github.com/dzhordano/go-posts/pkg/hasher"
	"github.com/dzhordano/go-posts/pkg/otp"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type (
	Users interface {
		SignUP(ctx context.Context, input UserSignUpInput) error
		SignIN(ctx context.Context, input UserSignInInput) (Tokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
		Verify(ctx context.Context, userId uint, codeHash string) error

		GetAll(ctx context.Context) ([]domain.User, error)
		GetById(ctx context.Context, userId uint) (domain.User, error)
	}

	UserSignUpInput struct {
		Name     string `json:"name" binding:"required,min=2,max=64"`
		Email    string `json:"email" binding:"required,email,max=64"`
		Password string `json:"password" binding:"required,min=8,max=64"`
	}

	UserSignInInput struct {
		Email    string `json:"email" binding:"required,email,max=64"`
		Password string `json:"password" binding:"required,min=8,max=64"`
	}

	Admins interface {
		SignIN(ctx context.Context, input UserSignInInput) (Tokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)

		UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error
		DeleteUser(ctx context.Context, userId uint) error

		SuspendUser(ctx context.Context, userId uint) error
		SuspendPost(ctx context.Context, postId uint) error

		CensorComment(ctx context.Context, commId uint) error
		DeleteComment(ctx context.Context, commId uint) error
		// TODO: mark as MARK report, not just delete. fix this in case arent lazy
		DealReport(ctx context.Context, reportId uint) error
	}

	Posts interface {
		GetAll(ctx context.Context) ([]domain.Post, error)
		GetById(ctx context.Context, postId uint) (domain.Post, error)
		Update(ctx context.Context, input domain.UpdatePostInput, postId uint) error
		Delete(ctx context.Context, postId uint) error

		Create(ctx context.Context, input domain.Post, userId uint) error
		GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error)
		UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error
		DeleteUser(ctx context.Context, postId, userId uint) error
		Like(ctx context.Context, postId uint) error
		RemoveLike(ctx context.Context, postId uint) error
		// TODO: Move those to another REPORTS SERVICE
		GetAllReports(ctx context.Context) ([]domain.Report, error)
		Report(ctx context.Context, postId, userId uint) error
	}

	Comments interface {
		Create(ctx context.Context, input domain.Comment, postId uint) error
		GetComments(ctx context.Context, postId uint) ([]domain.Comment, error)

		GetUserComments(ctx context.Context, userId uint) ([]domain.Comment, error)
		GetUserPostComments(ctx context.Context, userId, postId uint) ([]domain.Comment, error)
		UpdateUser(ctx context.Context, input domain.UpdateCommentInput, commId, userId uint) error
		DeleteUser(ctx context.Context, commId, userId uint) error
	}

	VerificationEmailInput struct {
		Email            string
		Name             string
		VerificationCode string
	}

	Emails interface {
		SendUserVerificationEmail(VerificationEmailInput) error
	}
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Services struct {
	Users    Users
	Admins   Admins
	Posts    Posts
	Comments Comments
	Emails   Emails
}

// services dependencies
type Deps struct {
	Repos        *repository.Repos
	Hasher       hasher.PassworsHasher
	TokenManager auth.TokenManager
	EmailSender  email.Sender
	OtpGenerator otp.Generator

	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeLength int
}

func NewService(deps Deps) *Services {
	emailsService := NewEmailsService(deps.EmailSender)
	commentsService := NewCommentsService(deps.Repos.Comments)
	postsService := NewPostsService(deps.Repos.Posts, commentsService)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, postsService, commentsService, emailsService, deps.OtpGenerator, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeLength)
	adminsService := NewAdminsService(deps.Repos.Admins, deps.Hasher, deps.TokenManager, postsService, usersService, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users:    usersService,
		Admins:   adminsService,
		Posts:    postsService,
		Comments: commentsService,
		Emails:   emailsService,
	}
}
