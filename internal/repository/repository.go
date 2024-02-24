package repository

import (
	"context"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	users_table  = "users"
	admins_table = "admins"
)

type (
	Users interface {
		Create(ctx context.Context, user domain.User) error
		GetById(ctx context.Context, userId int) (domain.User, error)
		GetByCredentials(ctx context.Context, input domain.UserSignInInput) (domain.User, error)
		CreateSession(ctx context.Context, userId int, session domain.Session) error
		GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	}

	Admins interface {
		GetById(ctx context.Context, userId int) (domain.User, error)
		GetByCredentials(ctx context.Context, input domain.UserSignInInput) (domain.User, error)
	}

	Posts interface {
		Create(ctx context.Context, title, description string) error
		GetAll(ctx context.Context) ([]domain.Post, error)
		GetById(ctx context.Context, postId int) (domain.Post, error)
		Update(ctx context.Context, input domain.UpdatePostInput) (domain.Post, error)
		Delete(ctx context.Context) error
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
	}
}
