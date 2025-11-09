package usecase

import (
	"time"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("foodcourt-secret-key-2025")

type AuthUseCase interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authUseCase struct {
	adminRepo repository.AdminRepository
}

func NewAuthUseCase(adminRepo repository.AdminRepository) AuthUseCase {
	return &authUseCase{adminRepo: adminRepo}
}

func (u *authUseCase) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	admin, err := u.adminRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,

		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	resp := &dto.LoginResponse{
		Token:   tokenString,
		Message: "Login successful",
	}
	resp.Admin.ID = admin.ID
	resp.Admin.Username = admin.Username
	resp.Admin.FullName = admin.FullName

	return resp, nil
}
