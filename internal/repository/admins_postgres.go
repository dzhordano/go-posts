package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminsRepo struct {
	db *pgxpool.Pool
}

func NewAdminsRepo(db *pgxpool.Pool) *AdminsRepo {
	return &AdminsRepo{
		db: db,
	}
}

func (r *AdminsRepo) GetById(ctx context.Context, userId uint) (domain.User, error) {
	panic("todo")
}

func (r *AdminsRepo) GetByCredentials(ctx context.Context, input domain.UserSignInInput) (domain.User, error) {
	query := fmt.Sprintf("SELECT id, name, email, password, registered, lastonline FROM %s WHERE email = $1 AND password = $2", admins_table)
	rows, err := r.db.Query(ctx, query, input.Email, input.Password)
	if err != nil {
		return domain.User{}, err
	}

	admin, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByNameLax[domain.User])
	if err != nil {
		fmt.Println(err)
		return domain.User{}, err
	}

	return *admin, nil
}

func (r *AdminsRepo) CreateSession(ctx context.Context, adminId uint, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET session.rtoken = $1, session.expiresat = $2, lastonline = $3 WHERE id = $4", admins_table)

	_, err := r.db.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, time.Now(), adminId)

	return err
}

func (r *AdminsRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE (session).rtoken = $1`, admins_table)

	row := r.db.QueryRow(ctx, query, refreshToken)

	var admin domain.User

	err := row.Scan(&admin)
	if err != nil {
		return domain.User{}, err
	}

	return admin, nil
}
