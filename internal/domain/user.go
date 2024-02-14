package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Verification struct {
		Code       string `json:"verification_code"`
		IsVerified bool   `json:"verificatio_verified"`
	} `json:"verification"`
	Suspended    bool      `json:"suspended"`
	RegisteredAt time.Time `json:"registered_at"`
	LastOnline   time.Time `json:"last_online"`
}
