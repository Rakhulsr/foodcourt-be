package utils

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func SetFlash(c *gin.Context, status string, message string) {

	encodedMsg := url.QueryEscape(message)

	c.SetCookie("flash_message", encodedMsg, 5, "/", "", false, true)
	c.SetCookie("flash_type", status, 5, "/", "", false, true)
}
