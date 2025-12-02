package repository

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

type WhatsAppLogRepository interface {
	Create(log *model.WhatsAppLog) error
	FindByOrderID(orderID uint) ([]model.WhatsAppLog, error)
	FindAll(page int, limit int) ([]model.WhatsAppLog, int64, error)
}

type whatsAppLogRepository struct {
	db *gorm.DB
}

func NewWhatsAppLogRepository(db *gorm.DB) WhatsAppLogRepository {
	return &whatsAppLogRepository{db: db}
}

func (r *whatsAppLogRepository) Create(log *model.WhatsAppLog) error {
	return r.db.Create(log).Error
}

func (r *whatsAppLogRepository) FindByOrderID(orderID uint) ([]model.WhatsAppLog, error) {
	var logs []model.WhatsAppLog
	err := r.db.Where("order_id = ?", orderID).Find(&logs).Error
	return logs, err
}

func (r *whatsAppLogRepository) FindAll(page int, limit int) ([]model.WhatsAppLog, int64, error) {
	var logs []model.WhatsAppLog
	var total int64

	offset := (page - 1) * limit

	r.db.Model(&model.WhatsAppLog{}).Count(&total)

	err := r.db.
		Preload("Order").
		Order("sent_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}
