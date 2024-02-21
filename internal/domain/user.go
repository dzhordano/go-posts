package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	Password     string    `json:"password"  binding:"required"`
	Verification `json:"verification"`
	Suspended    bool      `json:"suspended"`
	RegisteredAt time.Time `json:"registered_at"`
	LastOnline   time.Time `json:"last_online"`
}

type Verification struct {
	Code       string `json:"verification_code"`
	IsVerified bool   `json:"verificatio_verified"`
}

type UserSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserSignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
