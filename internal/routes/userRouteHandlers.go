package routes

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func setErrorResponse(ctx *gin.Context, erroCode int, message string) {
	ctx.JSON(erroCode, gin.H{
		"status":  "error",
		"message": message,
	})
}

// route handler for creating a new user
func CreateUser(ctx *gin.Context) {

	var user UserReq
	err := ctx.BindJSON(&user)
	if err != nil {
		setErrorResponse(ctx, 400, "Invalid request")
		return
	}
	// encrypt the password before saving it in DB using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			setErrorResponse(ctx, 400, "Password too long. Please provide a password with less than 72 characters")
			return
		}
		setErrorResponse(ctx, 500, "Failed to create User. Try again later!")
		return
	}
	user.Password = string(hashedPassword)
	// save the user in DB
	ctag, err := initializers.DB.Exec(context.Background(), "INSERT INTO passbook_app.users (username, email, password_hash,created_at,updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.Username,
		user.Email,
		user.Password,
		time.Now().UTC(), time.Now().UTC())
	log.Println(ctag)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			setErrorResponse(ctx, 400, "Username or Email already exists. Please provide a unique username and email.")
			return
		}
		setErrorResponse(ctx, 500, "Failed to create User. Try again later!")
		return
	}
	// get the generated userid
	var user_id string
	initializers.DB.QueryRow(context.Background(), "SELECT user_id FROM passbook_app.users WHERE username=$1", user.Username).Scan(&user_id)
	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"user_id":  user_id,
		},
		"meta": nil,
	})
}

// route handler for logging in a user
func LoginUser(ctx *gin.Context) {
	var userReq UserReq
	err := ctx.BindJSON(&userReq)
	if err != nil {
		setErrorResponse(ctx, 400, "Invalid request")
		return
	}
	rows, _ := initializers.DB.Query(context.Background(), "SELECT * FROM passbook_app.users WHERE username=$1", userReq.Username)
	// scan row into user struct
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		setErrorResponse(ctx, 401, "Invalid username or password")
		log.Println(err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userReq.Password))
	if err != nil {
		setErrorResponse(ctx, 401, "Invalid username or password")
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged in successfully",
		"data": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"user_id":  user.UserID,
		},
	})
}
