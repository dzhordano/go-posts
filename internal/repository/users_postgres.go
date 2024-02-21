package repository

import (
	"context"
	"fmt"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) error {
	query := fmt.Sprintf("INSERT INTO %s (name, email, password, registered, lastonline, verification.code, verification.verified) VALUES ($1, $2, $3, $4, $5, $6, $7)", users_table)

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password, user.RegisteredAt, user.LastOnline, user.Verification.Code, user.Verification.IsVerified)

	fmt.Println(err)
	return err
}

func (r *UsersRepo) GetById(ctx context.Context, userId uint) (*domain.User, error) {
	panic("todo")
}
