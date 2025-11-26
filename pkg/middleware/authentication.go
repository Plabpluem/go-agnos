package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *gin.Context) {
	accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	secretKey := "supersecret"

	token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Set("payload", claims)
	c.Next()
}
