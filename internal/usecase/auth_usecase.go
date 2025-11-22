package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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
		return nil, errors.New("username atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	secretKey := os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
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
