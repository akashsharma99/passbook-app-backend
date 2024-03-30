package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/akashsharma99/passbook-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth User middleware to check if the user is authenticated
func AuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check for the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized request",
			})
			return
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid Authorization header",
			})
			return
		}
		jwtToken := authHeaderParts[1]
		claims, err := ValidateToken(jwtToken, os.Getenv("ACCESS_SECRET"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid token",
			})
			return
		}
		ctx.Set("userId", claims.UserID)
		ctx.Next()
	}
}

func ValidateToken(jwtToken string, secret string) (*types.UserTokenClaims, error) {
	claims := &types.UserTokenClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !token.Valid {
		log.Println("Invalid token")
		return nil, err
	}
	return claims, nil
}
