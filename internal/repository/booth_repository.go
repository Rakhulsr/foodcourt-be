package repository

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

type BoothRepository interface {
	Create(booth *model.Booth) error
	FindAll() ([]model.Booth, error)
	FindActive() ([]model.Booth, error)
	FindByID(id uint) (*model.Booth, error)
	FindByName(keyword string) ([]model.Booth, error)
	FindByExactName(name string) (*model.Booth, error)
	Update(booth *model.Booth) error
	Delete(id uint) error
}

type BoothRepositoryImpl struct {
	db *gorm.DB
}

func NewBoothRepository(db *gorm.DB) *BoothRepositoryImpl {
	return &BoothRepositoryImpl{db: db}
}

func (r *BoothRepositoryImpl) Create(booth *model.Booth) error {
	return r.db.Create(booth).Error
}

func (r *BoothRepositoryImpl) FindAll() ([]model.Booth, error) {
	var booths []model.Booth
	err := r.db.Find(&booths).Error
	return booths, err
}

func (r *BoothRepositoryImpl) FindByID(id uint) (*model.Booth, error) {
	var booth model.Booth
	err := r.db.First(&booth, id).Error
	if err != nil {
		return nil, err
	}
	return &booth, nil
}

func (r *BoothRepositoryImpl) FindByExactName(name string) (*model.Booth, error) {
	var booth model.Booth
	err := r.db.Where("LOWER(name) = LOWER(?)", name).First(&booth).Error
	return &booth, err
}

func (r *BoothRepositoryImpl) FindByName(keyword string) ([]model.Booth, error) {
	var booths []model.Booth

	err := r.db.Preload("Menus").Where("LOWER(name) LIKE ?", "%"+keyword+"%").Find(&booths).Error

	return booths, err
}

func (r *BoothRepositoryImpl) FindActive() ([]model.Booth, error) {
	var booths []model.Booth
	err := r.db.Where("is_active = ?", true).Find(&booths).Error
	return booths, err
}

func (r *BoothRepositoryImpl) Update(booth *model.Booth) error {
	return r.db.Save(booth).Error
}

func (r *BoothRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&model.Booth{}, id).Error
}
