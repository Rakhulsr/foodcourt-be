package repository

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByCode(code string) (*model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return tx.Create(&order.Items).Error
	})
}

func (r *orderRepository) FindByCode(code string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("order_code = ?", code).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
