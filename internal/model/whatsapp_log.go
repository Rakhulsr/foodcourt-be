package model

import "time"

type WhatsAppLog struct {
	ID          uint      `gorm:"primaryKey"`
	OrderID     *uint     `gorm:"index"`
	BoothID     *uint     `gorm:"index"`
	MessageType string    `gorm:"size:50"`
	Status      string    `gorm:"size:20"`
	Response    string    `gorm:"type:text"`
	SentAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
