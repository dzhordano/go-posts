package repository

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	users_table    = "users"
	admins_table   = "admins"
	posts_table    = "posts"
	comments_table = "comments"
	reports_table  = "reports"

	users_posts = "users_posts"
)

type (
	Users interface {
		Create(ctx context.Context, user domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		GetById(ctx context.Context, userId uint) (domain.User, error)
		GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
		CreateSession(ctx context.Context, userId uint, session domain.Session) error
		GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
		Verify(ctx context.Context, userId uint, codeHash string) error
	}

	Admins interface {
		GetById(ctx context.Context, userId uint) (domain.User, error)
		GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
		CreateSession(ctx context.Context, adminId uint, session domain.Session) error
		GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)

		UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error
		DeleteUser(ctx context.Context, userId uint) error
		SuspendUser(ctx context.Context, userId uint) error
		SuspendPost(ctx context.Context, postId uint) error
		CensorComment(ctx context.Context, commId uint) error
		DeleteComment(ctx context.Context, commId uint) error
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
		Report(ctx context.Context, postId, userId uint) error
		GetAllReports(ctx context.Context) ([]domain.Report, error)
	}

	Comments interface {
		Create(ctx context.Context, input domain.Comment, postId uint) error
		GetComments(ctx context.Context, postId uint) ([]domain.Comment, error)

		UpdateUser(ctx context.Context, input domain.UpdateCommentInput, commId, userId uint) error
		DeleteUser(ctx context.Context, commId, userId uint) error
		GetUserComments(ctx context.Context, userId uint) ([]domain.Comment, error)
		GetUserPostComments(ctx context.Context, userId, postId uint) ([]domain.Comment, error)
	}
)

type Repos struct {
	Users    Users
	Admins   Admins
	Posts    Posts
	Comments Comments
}

func NewRepos(db *pgxpool.Pool) *Repos {
	return &Repos{
		Users:    NewUsersRepo(db),
		Admins:   NewAdminsRepo(db),
		Posts:    NewPostsRepo(db),
		Comments: NewCommentsRepo(db),
	}
}
