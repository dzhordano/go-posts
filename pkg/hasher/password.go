package hasher

import (
	"crypto/sha256"
	"fmt"
)

type PassworsHasher interface {
	GeneratePasswordHash(password string) (string, error)
}

type SHA256Hasher struct {
	salt string
}

func NewSHA256Hasher(salt string) *SHA256Hasher {
	return &SHA256Hasher{
		salt: salt,
	}
}

func (s *SHA256Hasher) GeneratePasswordHash(password string) (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt))), nil
}
