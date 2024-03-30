package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserTokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}
type User struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
