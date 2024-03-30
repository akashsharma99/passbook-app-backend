package routes

import (
	"context"
	"log"
	"time"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type User struct {
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetUser(ctx *gin.Context) {
	// take the user id from the auth middleware
	userID := ctx.MustGet("userId").(string)
	// get the user from the DB
	rows, _ := initializers.DB.Query(context.Background(), "SELECT * FROM passbook_app.users WHERE user_id=$1", userID)
	// scan row into user struct
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if err == pgx.ErrNoRows {
			setErrorResponse(ctx, 404, "User not found")
			return
		}
		setErrorResponse(ctx, 500, "Internal server error")
		log.Println(err)
		return
	}
	// return user in response body
	ctx.JSON(200, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"user": map[string]string{
				"username": user.Username,
				"email":    user.Email,
			},
		},
		"meta": nil,
	})
}
