package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const tokenDuration = 24 * time.Hour

func secret() ([]byte, error) {
	value := os.Getenv("JWT_SECRET")
	if value == "" {
		return nil, fmt.Errorf("JWT_SECRET is not configured")
	}
	return []byte(value), nil
}

func Issue(accountID int) (string, error) {
	secret, err := secret()
	if err != nil {
		return "", err
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(accountID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
	})
	return token.SignedString(secret)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, err := secret()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing bearer token"})
			return
		}
		token, err := jwt.ParseWithClaims(strings.TrimPrefix(header, "Bearer "), &jwt.RegisteredClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return secret, nil
			})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token claims"})
			return
		}
		accountID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token subject"})
			return
		}
		c.Set("accountID", accountID)
		c.Next()
	}
}
