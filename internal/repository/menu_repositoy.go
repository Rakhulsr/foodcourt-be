package repository

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(menu *model.Menu) error
	FindAll() ([]model.Menu, error)
	FindByID(id uint) (*model.Menu, error)
	FindByBoothID(boothID uint) ([]model.Menu, error)

	FindActive() ([]model.Menu, error)
	FindActiveByBoothID(boothID uint) ([]model.Menu, error)
	FindByName(keyword string) ([]model.Menu, error)
	FindByCategory(category string) ([]model.Menu, error)

	Update(menu *model.Menu) error
	Delete(id uint) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByBoothID(boothID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("booth_id = ?", boothID).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindActive() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Booth").Where("is_available = ?", true).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Preload("Booth").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) FindActiveByBoothID(boothID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Booth").Where("booth_id = ? AND is_available = ?", boothID, true).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByName(keyword string) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Booth").Where("LOWER(name) LIKE ?", "%"+keyword+"%").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindByCategory(category string) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Preload("Booth").Where("category = ?", category).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}
