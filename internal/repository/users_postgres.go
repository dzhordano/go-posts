package repository

import (
	"context"
	"fmt"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5"
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

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password, user.RegisteredAt, user.LastOnline, user.Verification.Code, user.Verification.Verified)

	fmt.Println(err)
	return err
}

func (r *UsersRepo) GetById(ctx context.Context, userId uint) (domain.User, error) {
	panic("todo")
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, input domain.UserSignInInput) (domain.User, error) {
	query := fmt.Sprintf("SELECT id, uid, name, email, password, (verification).code, (verification).verified, registered, lastonline FROM %s WHERE email = $1 AND password = $2", users_table)
	rows, err := r.db.Query(ctx, query, input.Email, input.Password)
	if err != nil {
		return domain.User{}, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.User])
	if err != nil {
		return domain.User{}, err
	}
	fmt.Println(*user)

	return *user, nil
}
