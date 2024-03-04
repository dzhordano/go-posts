package repository

import (
	"context"
	"fmt"

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
