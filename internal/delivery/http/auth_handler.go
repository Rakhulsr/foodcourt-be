package http

import (
	"net/http"
	"os"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) ShowLoginForm(c *gin.Context) {

	_, err := c.Cookie("admin_token")
	if err == nil {
		c.Redirect(http.StatusFound, "/api/admin/booths")
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"csrf_token": c.GetString("csrf_token"),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error":      "Username dan Password wajib diisi",
			"csrf_token": c.GetString("csrf_token"),
		})
		return
	}

	resp, err := h.authUC.Login(req)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error":      "Username atau Password salah",
			"csrf_token": c.GetString("csrf_token"),
		})
		return
	}

	isSecure := os.Getenv("GIN_MODE") == "release"
	c.SetCookie("admin_token", resp.Token, 3600*24, "/", "", isSecure, true)

	c.Redirect(http.StatusFound, "/api/admin/booths")
}

func (h *AuthHandler) Logout(c *gin.Context) {

	c.SetCookie("admin_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/auth/login")
}
