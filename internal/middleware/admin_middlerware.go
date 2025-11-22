package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		cookie, err := c.Cookie("admin_token")
		if err == nil {
			tokenStr = cookie
		}

		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			handleUnauthorized(c)
			return
		}

		jwtSecret := []byte(os.Getenv("SECRET_KEY"))
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			handleUnauthorized(c)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("admin_id", claims["admin_id"])
		c.Next()
	}
}

func handleUnauthorized(c *gin.Context) {

	if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	c.Abort()
}
