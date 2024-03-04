package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dzhordano/go-posts/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostsRepo struct {
	db *pgxpool.Pool
}

func NewPostsRepo(db *pgxpool.Pool) *PostsRepo {
	return &PostsRepo{
		db: db,
	}
}

func (r *PostsRepo) Create(ctx context.Context, input domain.Post, userId uint) error {
	// TODO: reimplement with r.db.BeginTx()
	qPostsTable := fmt.Sprintf("INSERT INTO %s (title, description, created, updated) VALUES ($1, $2, $3, $4) RETURNING id", posts_table)

	var postId int
	row := r.db.QueryRow(ctx, qPostsTable, input.Title, input.Description, time.Now(), time.Now())

	err := row.Scan(&postId)
	if err != nil {
		return err
	}

	qUsersPosts := fmt.Sprintf("INSERT INTO %s (post_id, user_id) VALUES ($1, $2)", users_posts)

	_, err = r.db.Exec(ctx, qUsersPosts, postId, userId)

	return err
}

func (r *PostsRepo) GetAll(ctx context.Context) ([]domain.Post, error) {
	var posts []domain.Post

	query := fmt.Sprintf("SELECT id, title, description, suspended, created, updated, likes, watched FROM %s", posts_table)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []domain.Post{}, err
	}

	posts, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Post])
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostsRepo) GetById(ctx context.Context, postId uint) (domain.Post, error) {
	var post domain.Post

	query := fmt.Sprintf("SELECT id, title, description, suspended, created, updated, likes, watched FROM %s WHERE id = $1", posts_table)

	row := r.db.QueryRow(ctx, query, postId)

	// TODO: do i need to set session values?
	err := row.Scan(&post.ID, &post.Title, &post.Description, &post.Suspended, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Watched)
	if err != nil {
		return domain.Post{}, err
	}

	return post, nil
}

func (r *PostsRepo) Update(ctx context.Context, input domain.UpdatePostInput) (domain.Post, error) {
	panic("TODO")
}

func (r *PostsRepo) Delete(ctx context.Context) error {
	panic("TODO")
}

func (r *PostsRepo) GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error) {
	var posts []domain.Post

	query := fmt.Sprintf("SELECT p.id, p.title, p.description, p.suspended, p.created, p.updated, p.likes, p.watched FROM %s p JOIN %s up ON p.id = up.post_id WHERE up.user_id = $1", posts_table, users_posts)

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return []domain.Post{}, err
	}

	posts, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Post])
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostsRepo) GetByIdUser(ctx context.Context, postId, userId uint) (domain.Post, error) {
	var post domain.Post

	query := fmt.Sprintf("SELECT p.id, p.title, p.description, p.suspended, p.created, p.updated, p.likes, p.watched FROM %s p JOIN %s up ON p.id = up.post_id WHERE p.id = $1 AND up.user_id = $2", posts_table, users_posts)

	row := r.db.QueryRow(ctx, query, postId, userId)

	// TODO: do i need to set session values?
	err := row.Scan(&post.ID, &post.Title, &post.Description, &post.Suspended, &post.CreatedAt, &post.UpdatedAt, &post.Likes, &post.Watched)
	if err != nil {
		return domain.Post{}, err
	}

	return post, nil
}