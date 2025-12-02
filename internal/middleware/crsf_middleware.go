package middleware

import (
	"net/http"
	"strings"

	"github.com/Rakhulsr/foodcourt/utils"
	"github.com/gin-gonic/gin"
)

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/webhooks/") {
			c.Next()
			return
		}

		csrfToken, err := c.Cookie("csrf_token")

		if err != nil || csrfToken == "" {
			token := utils.RandomString(32)

			c.SetCookie("csrf_token", token, 3600*24, "/", "", false, false)
			csrfToken = token
		}

		c.Set("csrf_token", csrfToken)

		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" || c.Request.Method == "PATCH" {

			inputToken := c.PostForm("csrf_token")

			if inputToken == "" {
				inputToken = c.GetHeader("X-CSRF-Token")
			}

			if inputToken == "" || inputToken != csrfToken {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch. Refresh page."})
				return
			}
		}

		c.Next()
	}
}
