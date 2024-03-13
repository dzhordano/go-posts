package repository

import (
	"context"
	"fmt"
	"strings"
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
	query := fmt.Sprintf("SELECT id, name, email, password, registered, lastonline FROM %s WHERE id = $1", admins_table)
	row := r.db.QueryRow(ctx, query, userId)

	var admin domain.User

	err := row.Scan(&admin)
	if err != nil {
		return domain.User{}, err
	}

	return admin, nil
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

func (r *AdminsRepo) UpdateUser(ctx context.Context, input domain.UpdateUserInput, userId uint) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	if input.Verification != nil {
		setValues = append(setValues, fmt.Sprintf("(verification).code=$%d", argId))
		args = append(args, *input.Password)
		argId++

		setValues = append(setValues, fmt.Sprintf("(verification).verified=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	if input.Suspended != nil {
		setValues = append(setValues, fmt.Sprintf("suspended=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", users_table, setQuery, argId)

	args = append(args, userId)

	_, err := r.db.Exec(ctx, query, args...)

	return err
}

func (r *AdminsRepo) DeleteUser(ctx context.Context, userId uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", users_table)

	_, err := r.db.Exec(ctx, query, userId)

	return err
}

func (r *AdminsRepo) SuspendUser(ctx context.Context, userId uint) error {
	query := fmt.Sprintf("UPDATE %s SET suspended = TRUE WHERE id = $1", users_table)

	_, err := r.db.Exec(ctx, query, userId)

	return err
}

func (r *AdminsRepo) SuspendPost(ctx context.Context, postId uint) error {
	query := fmt.Sprintf("UPDATE %s SET suspended = TRUE WHERE id = $1", posts_table)

	_, err := r.db.Exec(ctx, query, postId)

	return err
}

func (r *AdminsRepo) CensorComment(ctx context.Context, commId uint) error {
	query := fmt.Sprintf("UPDATE %s SET censored = TRUE WHERE post_id = $1 AND id = $2", comments_table)

	_, err := r.db.Exec(ctx, query, commId)

	return err
}

func (r *AdminsRepo) DeleteComment(ctx context.Context, commId uint) error {
	qCommentsTable := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING post_id", comments_table)

	var postId int
	row := r.db.QueryRow(ctx, qCommentsTable, commId)

	if err := row.Scan(&postId); err != nil {
		return err
	}

	qPostsTable := fmt.Sprintf("UPDATE %s SET comments = comments - 1 WHERE id = $1", posts_table)

	_, err := r.db.Exec(ctx, qPostsTable, postId)

	return err
}
