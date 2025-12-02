package repository

import (
	"time"

	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByCode(code string) (*model.Order, error)
	UpdateInvoice(orderCode string, invoiceID string, invoiceURL string) error

	FindAll(page int, limit int, status string) ([]model.Order, int64, error)
	UpdatePaymentStatus(orderCode string, status string) error
	UpdateOrderStatus(orderCode string, status string) error

	GetTotalIncomeToday() (int, error)
	CountOrdersToday() (int64, error)
	FindOrdersToday() ([]model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Omit("Items.ID").Create(order).Error
}

func (r *orderRepository) FindByCode(code string) (*model.Order, error) {
	var order model.Order
	err := r.db.
		Preload("Items").
		Preload("Items.Menu").
		Preload("Items.Booth").
		Where("order_code = ?", code).
		First(&order).Error

	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateInvoice(orderCode string, invoiceID string, invoiceURL string) error {

	return r.db.Model(&model.Order{}).
		Where("order_code = ?", orderCode).
		Updates(map[string]interface{}{
			"xendit_invoice_id": invoiceID,
			"invoice_url":       invoiceURL,
		}).Error
}

func (r *orderRepository) FindAll(page int, limit int, status string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * limit

	query := r.db.Model(&model.Order{})

	if status != "" {

		query = query.Where("order_status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Items").
		Preload("Items.Menu").
		Preload("Items.Booth").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, total, err
}
func (r *orderRepository) UpdatePaymentStatus(orderCode string, status string) error {
	return r.db.Model(&model.Order{}).
		Where("order_code = ?", orderCode).
		Update("payment_status", status).Error
}

func (r *orderRepository) UpdateOrderStatus(orderCode string, status string) error {
	return r.db.Model(&model.Order{}).
		Where("order_code = ?", orderCode).
		Update("order_status", status).Error
}

func (r *orderRepository) GetTotalIncomeToday() (int, error) {
	var total int

	today := time.Now().Truncate(24 * time.Hour)

	err := r.db.Model(&model.Order{}).
		Where("order_status = ? AND created_at >= ?", "completed", today).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error

	return total, err
}

func (r *orderRepository) CountOrdersToday() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)

	err := r.db.Model(&model.Order{}).
		Where("created_at >= ?", today).
		Count(&count).Error

	return count, err
}

func (r *orderRepository) FindOrdersToday() ([]model.Order, error) {
	var orders []model.Order
	today := time.Now().Truncate(24 * time.Hour)

	err := r.db.
		Preload("Items").
		Preload("Items.Menu").
		Preload("Items.Booth").
		Where("created_at >= ?", today).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}
