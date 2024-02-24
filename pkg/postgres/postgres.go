package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
	MaxAtts  int
}

// TODO: implement those or what? damn...
type Client interface {
	Exec(ctx context.Context, sql string, argument ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, dc DBConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dc.Username, dc.Password, dc.Host, dc.Port, dc.Database, dc.SSLMode)

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	return dbpool, nil
}
