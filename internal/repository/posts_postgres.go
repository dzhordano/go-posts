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
	// this inserts into posts_table the post we create and also pulls name of author from db
	qPostsTable := fmt.Sprintf("INSERT INTO %s (title, description, author, created, updated) VALUES ($1, $2, $3, $4, $5) RETURNING id", posts_table)

	var postId int
	row := r.db.QueryRow(ctx, qPostsTable, input.Title, input.Description, input.Author, time.Now(), time.Now())

	if err := row.Scan(&postId); err != nil {
		return err
	}

	qUsersPosts := fmt.Sprintf("INSERT INTO %s (post_id, user_id) VALUES ($1, $2)", users_posts)

	_, err := r.db.Exec(ctx, qUsersPosts, postId, userId)

	return err
}

func (r *PostsRepo) GetAll(ctx context.Context) ([]domain.Post, error) {
	query := fmt.Sprintf("UPDATE %s SET watched = watched + 1 RETURNING *", posts_table)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []domain.Post{}, err
	}

	var posts []domain.Post
	posts, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Post])
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostsRepo) GetById(ctx context.Context, postId uint) (domain.Post, error) {
	var post domain.Post

	query := fmt.Sprintf("UPDATE %s SET watched = watched + 1 WHERE id = $1 RETURNING *", posts_table)

	row := r.db.QueryRow(ctx, query, postId)

	// TODO: do i need to set session values?
	err := row.Scan(&post.ID, &post.Title, &post.Description, &post.Author, &post.Comments, &post.Suspended, &post.Created, &post.Updated, &post.Likes, &post.Watched)
	if err != nil {
		return domain.Post{}, err
	}

	return post, nil
}

func (r *PostsRepo) Update(ctx context.Context, input domain.UpdatePostInput, postId uint) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", posts_table, setQuery, argId)

	args = append(args, postId)

	_, err := r.db.Exec(ctx, query, args...)

	return err
}

func (r *PostsRepo) Delete(ctx context.Context, postId uint) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", posts_table)

	_, err := r.db.Exec(ctx, query, postId)

	return err
}

func (r *PostsRepo) GetAllUser(ctx context.Context, userId uint) ([]domain.Post, error) {
	var posts []domain.Post

	query := fmt.Sprintf("SELECT p.id, p.title, p.description, p.author, p.comments, p.suspended, p.created, p.updated, p.likes, p.watched FROM %s p INNER JOIN %s up ON p.id = up.post_id WHERE up.user_id = $1", posts_table, users_posts)

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

func (r *PostsRepo) UpdateUser(ctx context.Context, input domain.UpdatePostInput, postId, userId uint) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	// CHECK IF VALUES ARE NIL AND IF NOT -> ADD THEM TO QUERY STRING AS setQuery
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s p SET %s FROM %s up WHERE p.id = up.post_id AND up.post_id = $%d AND up.user_id = $%d", posts_table, setQuery, users_posts, argId, argId+1)

	args = append(args, postId, userId)

	_, err := r.db.Exec(ctx, query, args...)

	return err
}

func (r *PostsRepo) DeleteUser(ctx context.Context, postId, userId uint) error {
	query := fmt.Sprintf("DELETE FROM %s p USING %s up WHERE p.id = up.post_id AND up.post_id = $1 AND up.user_id = $2", posts_table, users_posts)

	_, err := r.db.Exec(ctx, query, postId, userId)

	return err
}

func (r *PostsRepo) Like(ctx context.Context, postId uint) error {
	query := fmt.Sprintf("UPDATE %s SET likes = likes + 1 WHERE id = $1", posts_table)

	_, err := r.db.Exec(ctx, query, postId)

	return err
}

func (r *PostsRepo) RemoveLike(ctx context.Context, postId uint) error {
	query := fmt.Sprintf("UPDATE %s SET likes = likes - 1 WHERE id = $1", posts_table)

	_, err := r.db.Exec(ctx, query, postId)

	return err
}
