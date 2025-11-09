// internal/repository/whatsapp_log_repository.go
package repository

import (
	"github.com/Rakhulsr/foodcourt/internal/model"
	"gorm.io/gorm"
)

// === INTERFACE ===
type WhatsAppLogRepository interface {
	Create(log *model.WhatsAppLog) error
	FindByOrderID(orderID uint) ([]model.WhatsAppLog, error)
}

// === IMPLEMENTATION ===
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
