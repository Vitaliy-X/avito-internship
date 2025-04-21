package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "moderator" && token != "employee" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Invalid role"})
			return
		}

		c.Set("role", token)
		c.Next()
	}
}
