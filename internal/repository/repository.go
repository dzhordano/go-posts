package repository

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	users_table  = "users"
	admins_table = "admins"
	posts_table  = "posts"
	users_posts  = "users_posts"
)

type (
	Users interface {
		Create(ctx context.Context, user domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		GetById(ctx context.Context, userId uint) (domain.User, error)
		GetByCredentials(ctx context.Context, input domain.UserSignInInput) (domain.User, error)
		CreateSession(ctx context.Context, userId uint, session domain.Session) error
		GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	}

	Admins interface {
		GetById(ctx context.Context, userId uint) (domain.User, error)
		GetByCredentials(ctx context.Context, input domain.UserSignInInput) (domain.User, error)
		CreateSession(ctx context.Context, adminId uint, session domain.Session) error
		GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
		UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error
		DeleteUser(ctx context.Context, userId uint) error
	}

	Posts interface {
		Create(ctx context.Context, input domain.Post, userId uint) error
		GetAll(ctx context.Context) ([]domain.Post, error)
		GetById(ctx context.Context, postId uint) (domain.Post, error)
		Update(ctx context.Context, input domain.UpdatePostInput, postId uint) error
		Delete(ctx context.Context, postId uint) error

		GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error)
		GetByIdUser(ctx context.Context, postId, userId uint) (domain.Post, error)
		UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error
		DeleteUser(ctx context.Context, postId, userId uint) error
	}
)

type Repos struct {
	Users  Users
	Admins Admins
	Posts  Posts
}

func NewRepos(db *pgxpool.Pool) *Repos {
	return &Repos{
		Users:  NewUsersRepo(db),
		Admins: NewAdminsRepo(db),
		Posts:  NewPostsRepo(db),
	}
}
