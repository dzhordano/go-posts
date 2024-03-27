package domain

import (
	"errors"
	"time"
)

type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title" binding:"required,min=1"`
	Description string    `json:"description" binding:"required,min=1"`
	Author      string    `json:"author"`
	Suspended   bool      `json:"suspended"`
	Comments    int       `json:"comments"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Likes       uint      `json:"likes"`
	Watched     uint      `json:"watched"`
}

type Comment struct {
	ID       uint      `json:"id"`
	PostId   uint      `json:"post_id"`
	AuthorId uint      `json:"author_id"`
	Data     string    `json:"data" binding:"required,min=1"`
	Created  time.Time `json:"commented_at"`
	Updated  time.Time `json:"updated_at"`
	Censored bool      `json:"censored"`
}

type Report struct {
	ID         uint      `json:"id"`
	PostId     uint      `json:"post_id"`
	UserId     uint      `json:"user_id"`
	ReportedAt time.Time `json:"reported_at" db:"created"`
}

type UpdatePostInput struct {
	Title       *string
	Description *string
}

func (i UpdatePostInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update input has no values")
	}

	return nil
}

type UpdateCommentInput struct {
	Data string `json:"data" binding:"required,min=1"`
}
