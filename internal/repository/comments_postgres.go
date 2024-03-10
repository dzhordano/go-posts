package repository

import (
	"context"
	"fmt"

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
	panic("TODO")
}

func (r *CommentsRepo) GetComments(ctx context.Context, postId uint) ([]domain.Comment, error) {
	var comments []domain.Comment

	query := fmt.Sprintf("SELECT c.id, c.author, c.comment, c.created, c.updated, c.censored FROM %s c INNER JOIN %s pc ON c.id = pc.comment_id WHERE pc.post_id = $1", comments_table, posts_comments)

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
