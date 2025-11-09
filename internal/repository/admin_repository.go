package repository

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

type AdminRepository interface {
	FindByUsername(username string) (*model.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
