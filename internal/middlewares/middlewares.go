package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserTokenClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

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
		claims, err := ValidateToken(jwtToken)
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

func ValidateToken(jwtToken string) (*UserTokenClaims, error) {
	claims := &UserTokenClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil {
		if err == jwt.ErrTokenExpired {
			log.Println("Token expired")
		} else {
			log.Println(err)
		}
		return nil, err
	}
	if !token.Valid {
		log.Println("Invalid token")
		return nil, err
	}
	return claims, nil
}
