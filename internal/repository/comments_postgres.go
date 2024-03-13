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

func (r *CommentsRepo) UpdateUser(ctx context.Context, input domain.UpdateCommentInput, commId, userId uint) error {
	query := fmt.Sprintf("UPDATE %s SET data = $1, updated = $2 WHERE id = $3 AND author_id = $4", comments_table)

	_, err := r.db.Exec(ctx, query, input.Data, time.Now(), commId, userId)

	return err
}

func (r *CommentsRepo) DeleteUser(ctx context.Context, commId, userId uint) error {
	qCommentsTable := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND author_id = $2 RETURNING post_id", comments_table)

	var postId int
	row := r.db.QueryRow(ctx, qCommentsTable, commId, userId)

	if err := row.Scan(&postId); err != nil {
		return err
	}

	qPostsTable := fmt.Sprintf("UPDATE %s SET comments = comments - 1 WHERE id = $1", posts_table)

	_, err := r.db.Exec(ctx, qPostsTable, postId)

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
