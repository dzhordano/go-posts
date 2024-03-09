package domain

import "time"

type Session struct {
	RefreshToken string    `json:"refresh_token" db:"rtoken" `
	ExpiresAt    time.Time `json:"expires_at"`
}
