package middleware

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func FlashMessage() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("flash_message")
		if err == nil && cookie != "" {

			msg, _ := url.QueryUnescape(cookie)
			c.Set("FlashMessage", msg)

			cookieType, _ := c.Cookie("flash_type")
			if cookieType == "" {
				cookieType = "success"
			}
			c.Set("FlashType", cookieType)

			c.SetCookie("flash_message", "", -1, "/", "", false, true)
			c.SetCookie("flash_type", "", -1, "/", "", false, true)
		}
		c.Next()
	}
}
