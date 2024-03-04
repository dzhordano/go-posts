package domain

import (
	"time"
)

type Post struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"  binding:"required"`
	Suspended   bool      `json:"suspended"`
	CreatedAt   time.Time `json:"createdAt" db:"created"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated"`
	Likes       uint      `json:"likes"`
	Watched     uint      `json:"watched"`
}

type Comment struct {
	ID          uint   `json:"id"`
	UserName    string `json:"username"`
	Data        string `json:"data"`
	CommentedAt string `json:"commented_at"`
	UpdatedAt   string `json:"updated_at"`
	Censored    bool   `json:"censored"`
}

type UpdatePostInput struct {
	Title       *string
	Description *string
}
