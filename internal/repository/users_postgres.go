package repository

import (
	"context"
	"fmt"
	"time"

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

	return err
}

func (r *UsersRepo) GetById(ctx context.Context, userId int) (domain.User, error) {
	// FIXME: test this before moving on
	query := fmt.Sprintf("SELECT id, uid, name, email, password, (verification).code, (verification).verified, registered, lastonline FROM %s WHERE id = $1", users_table)
	rows, err := r.db.Query(ctx, query, userId)
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

	return *user, nil
}

func (r *UsersRepo) CreateSession(ctx context.Context, userId int, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET session.rtoken = $1, session.expiresat = $2, lastonline = $3 WHERE id = $4", users_table)

	_, err := r.db.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, time.Now(), userId)

	return err
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE (session).rtoken = $1`, users_table)
	rows, err := r.db.Query(ctx, query, refreshToken)
	if err != nil {
		return domain.User{}, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.User])
	if err != nil {
		return domain.User{}, err
	}

	return *user, nil
}
