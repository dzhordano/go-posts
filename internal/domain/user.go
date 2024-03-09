package domain

import (
	"errors"
	"time"
)

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password"  binding:"required"`
	Verification `json:"verification" db:"verification"`
	Session      `json:"session"`
	Suspended    bool      `json:"suspended"`
	RegisteredAt time.Time `json:"registered_at"  db:"registered"`
	LastOnline   time.Time `json:"last_online"  db:"lastonline"`
}

type Verification struct {
	Code     string `json:"verification_code"`
	Verified bool   `json:"verificatio_verified"`
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

type UpdateUserInput struct {
	Name     *string
	Password *string
	*Verification
	Suspended *bool
}

func (i UpdateUserInput) Validate() error {
	if i.Name == nil && i.Password == nil && i.Verification == nil && i.Suspended == nil {
		return errors.New("update input has no values")
	}

	return nil
}
