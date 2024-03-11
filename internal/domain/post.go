package domain

import (
	"errors"
	"time"
)

// FIXME: change createdAt -> created and updatedAt -> updated
type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"  binding:"required"`
	Author      string    `json:"author"`
	Suspended   bool      `json:"suspended"`
	Comments    int       `json:"comments"`
	CreatedAt   time.Time `json:"createdAt" db:"created"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated"`
	Likes       uint      `json:"likes"`
	Watched     uint      `json:"watched"`
}

type Comment struct {
	ID       uint `json:"id"`
	PostId   uint `json:"post_id"`
	AuthorId uint `json:"author_id"`
	// FIXME: change in sql table comment -> data
	Data     string    `json:"data" binding:"required" db:"comment"`
	Created  time.Time `json:"commented_at"`
	Updated  time.Time `json:"updated_at"`
	Censored bool      `json:"censored"`
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
