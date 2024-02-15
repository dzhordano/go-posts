package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Likes       uint      `json:"likes"`
	Watched     uint      `json:"watched"`
	Comment     Comment   `json:"comment"`
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
