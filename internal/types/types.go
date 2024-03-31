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
type Passbook struct {
	PassbookID    string    `json:"passbook_id"`
	UserID        string    `json:"user_id"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"account_number"`
	TotalBalance  float64   `json:"total_balance"`
	Nickname      string    `json:"nickname"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type Transaction struct {
	TransactionID   string    `json:"transaction_id"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	TransactionType string    `json:"transaction_type"`
	PartyName       string    `json:"party_name"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Tags            string    `json:"tags"`
	PassbookID      string    `json:"passbook_id"`
	UserID          string    `json:"user_id"`
}
