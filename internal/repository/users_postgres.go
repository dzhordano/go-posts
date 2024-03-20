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
	query := fmt.Sprintf("INSERT INTO %s (name, email, password, registered, lastonline, verification.code, verification.verified, session.rtoken, session.expiresat) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", users_table)

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.Password, user.RegisteredAt, user.LastOnline, user.Verification.Code, user.Verification.Verified, user.Session.RefreshToken, user.Session.ExpiresAt)

	return err
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	query := fmt.Sprintf("SELECT id, name, email, password, (verification).code, (verification).verified, suspended, registered, lastonline FROM %s WHERE email = $1 AND password = $2", users_table)

	row := r.db.QueryRow(ctx, query, email, password)

	var user domain.User

	// TODO: do i need to set session values?
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Verification.Code, &user.Verification.Verified, &user.Suspended, &user.RegisteredAt, &user.LastOnline)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) CreateSession(ctx context.Context, userId uint, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET session.rtoken = $1, session.expiresat = $2, lastonline = $3 WHERE id = $4", users_table)

	_, err := r.db.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, time.Now(), userId)

	return err
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE (session).rtoken = $1`, users_table)

	row := r.db.QueryRow(ctx, query, refreshToken)

	var user domain.User

	err := row.Scan(&user.ID)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
func (r *UsersRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	query := fmt.Sprintf("SELECT id, name, email, password, (verification).code, (verification).verified, (session).rtoken, (session).expiresat, suspended, registered, lastonline FROM %s", users_table)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []domain.User{}, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		fmt.Println(users)
		return []domain.User{}, err
	}

	return users, nil
}

func (r *UsersRepo) GetById(ctx context.Context, userId uint) (domain.User, error) {
	query := fmt.Sprintf("SELECT id, name, email, password, (verification).code, (verification).verified, suspended, registered, lastonline FROM %s WHERE id = $1", users_table)

	row := r.db.QueryRow(ctx, query, userId)

	var user domain.User

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Verification.Code, &user.Verification.Verified, &user.Suspended, &user.RegisteredAt, &user.LastOnline)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
