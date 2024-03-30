package routes

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/middlewares"
	"github.com/akashsharma99/passbook-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	// if any of the mandatory fields are empty, return an error
	if user.Username == "" || user.Email == "" || user.Password == "" {
		setErrorResponse(ctx, 400, "Please provide all the mandatory fields.")
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
			setErrorResponse(ctx, 409, "Username or Email already exists. Please provide a unique username and email.")
			return
		}
		setErrorResponse(ctx, 500, "Failed to create User. Try again later!")
		return
	}

	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data": map[string]interface{}{
			"user": map[string]string{
				"username": user.Username,
				"email":    user.Email,
			},
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
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[types.User])
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
	// generate access and refresh tokens
	access_token, refresh_token, err := generateTokens(user)
	if err != nil {
		setErrorResponse(ctx, 500, "Login failed. Try again later!")
		return
	}
	// return access token in response body while refresh token in httponly cookie
	ctx.SetCookie("refresh_token", refresh_token, 3600*24, "/", "", true, true)
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "User logged in successfully",
		"data": map[string]interface{}{
			"access_token": access_token,
			"user": map[string]string{
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}
func generateTokens(user types.User) (string, string, error) {
	time_now := time.Now()
	// generate signed access token
	accessClaims := types.UserTokenClaims{
		UserID: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time_now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time_now),
		},
	}
	access_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		log.Println("Failed to generate access token for user ", user.Username)
		return "", "", err
	}
	// generate signed refresh token
	refreshClaims := types.UserTokenClaims{
		UserID: user.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time_now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time_now),
		},
	}
	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		log.Println("Failed to generate refresh token for user ", user.Username)
		return "", "", err
	}
	// if token for user_id exists update it else insert token
	var exists bool
	dberr := initializers.DB.QueryRow(context.Background(), "SELECT true FROM passbook_app.tokens WHERE user_id=$1", user.UserID).Scan(&exists)
	if dberr != nil && dberr != pgx.ErrNoRows {
		log.Println(dberr)
		log.Println("Failed to check if token exists for user ", user.Username)
		return "", "", dberr
	}
	if exists {
		_, dberr = initializers.DB.Exec(context.Background(), "UPDATE passbook_app.tokens SET rtoken=$1, updated_at=$2 WHERE user_id=$3",
			refresh_token, time_now.UTC(), user.UserID)
	} else {
		_, dberr = initializers.DB.Exec(context.Background(), "INSERT INTO passbook_app.tokens (user_id, rtoken, created_at, updated_at) VALUES ($1, $2, $3, $4)",
			user.UserID, refresh_token, time_now.UTC(), time_now.UTC())
	}

	if dberr != nil {
		log.Println("Failed to save refresh token for user ", user.Username)
		return "", "", dberr
	}
	return access_token, refresh_token, nil
}
func RefreshToken(ctx *gin.Context) {
	refresh_token, err := ctx.Cookie("refresh_token")
	if err != nil {
		log.Println("Failed to get refresh token from cookie")
		setErrorResponse(ctx, 400, "Invalid request")
		return
	}
	claims, err := middlewares.ValidateToken(refresh_token, os.Getenv("REFRESH_SECRET"))
	if err != nil {
		setErrorResponse(ctx, 401, "Invalid Refresh token")
		return
	}
	// check if user exists and token is not revoked
	user_id := claims.UserID
	if !isValidUser(user_id) || isRevokedToken(refresh_token, user_id) {
		setErrorResponse(ctx, 401, "Invalid Refresh token")
		return
	}
	// generate new access and refresh tokens
	access_token, refresh_token, err := generateTokens(types.User{UserID: user_id})
	if err != nil {
		setErrorResponse(ctx, 500, "Failed to refresh token. Try again later!")
		return
	}
	// return access token in response body while refresh token in httponly cookie
	ctx.SetCookie("refresh_token", refresh_token, 3600*24, "/", "", true, true)
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Token refreshed successfully",
		"data": map[string]string{
			"access_token": access_token,
		},
	})
}
func isValidUser(user_id string) bool {
	var exists bool
	dberr := initializers.DB.QueryRow(context.Background(), "SELECT true FROM passbook_app.users WHERE user_id=$1", user_id).Scan(&exists)
	if dberr != nil && dberr != pgx.ErrNoRows {
		log.Println(dberr)
		log.Println("Failed to check if user exists for id ", user_id)
		return false
	}
	return exists
}
func isRevokedToken(refresh_token string, user_id string) bool {
	// if token not present in token table then it is a revoked token and should not be allowed to refresh
	var exists bool
	dberr := initializers.DB.QueryRow(context.Background(), "SELECT true FROM passbook_app.tokens WHERE user_id=$1 AND rtoken=$2", user_id, refresh_token).Scan(&exists)
	if dberr == pgx.ErrNoRows {
		log.Println("The incoming refresh token is a revoked token")
		return true
	} else if dberr != nil {
		log.Println("Failed to check token is revoked or not", dberr)
		return true
	}
	return !exists
}
