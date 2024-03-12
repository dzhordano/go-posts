package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentsRepo struct {
	db *pgxpool.Pool
}

func NewCommentsRepo(db *pgxpool.Pool) *CommentsRepo {
	return &CommentsRepo{
		db: db,
	}
}

func (r *CommentsRepo) Create(ctx context.Context, input domain.Comment, postId uint) error {
	// TODO: reimplement with tx.Begin()
	qCommsTable := fmt.Sprintf("INSERT INTO %s (post_id, author_id, data, created, updated) VALUES ($1, $2, $3, $4, $5)", comments_table)

	_, err := r.db.Exec(ctx, qCommsTable, postId, input.AuthorId, input.Data, time.Now(), time.Now())
	if err != nil {
		return err
	}

	qPostsCommCount := fmt.Sprintf("UPDATE %s SET comments = comments + 1 WHERE id = $1", posts_table)
	_, err = r.db.Exec(ctx, qPostsCommCount, postId)

	return err
}

func (r *CommentsRepo) GetComments(ctx context.Context, postId uint) ([]domain.Comment, error) {
	var comments []domain.Comment

	query := fmt.Sprintf("SELECT id, post_id, author_id, data, created, updated, censored FROM %s WHERE post_id = $1", comments_table)

	rows, err := r.db.Query(ctx, query, postId)
	if err != nil {
		return []domain.Comment{}, err
	}

	comments, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Comment])
	if err != nil {
		return []domain.Comment{}, err
	}

	return comments, nil
}

func (r *CommentsRepo) Delete(ctx context.Context, postId, commId uint) error {
	qCommentsTable := fmt.Sprintf("DELETE FROM %s WHERE post_id = $1 AND id = $2", comments_table)

	_, err := r.db.Exec(ctx, qCommentsTable, postId, commId)
	if err != nil {
		return err
	}

	qPostsTable := fmt.Sprintf("UPDATE %s SET comments = comments - 1 WHERE id = $1", posts_table)

	_, err = r.db.Exec(ctx, qPostsTable, postId)

	return err
}

func (r *CommentsRepo) GetUserComments(ctx context.Context, userId uint) ([]domain.Comment, error) {
	query := fmt.Sprintf("SELECT id, post_id, author_id, data, created, updated, censored FROM %s WHERE author_id = $1", comments_table)

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return []domain.Comment{}, err
	}

	var comments []domain.Comment
	comments, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Comment])
	if err != nil {
		return []domain.Comment{}, err
	}

	return comments, nil
}

func (r *CommentsRepo) GetUserPostComments(ctx context.Context, userId, postId uint) ([]domain.Comment, error) {
	query := fmt.Sprintf("SELECT id, post_id, author_id, data, created, updated, censored FROM %s WHERE author_id = $1 AND post_id = $2", comments_table)

	rows, err := r.db.Query(ctx, query, userId, postId)
	if err != nil {
		return []domain.Comment{}, err
	}

	var comments []domain.Comment
	comments, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Comment])
	if err != nil {
		return []domain.Comment{}, err
	}

	return comments, nil
}
