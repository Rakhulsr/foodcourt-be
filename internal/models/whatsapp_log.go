package models

import "time"

type WhatsAppLog struct {
	ID          uint `gorm:"primaryKey"`
	OrderID     *uint
	BoothID     *uint
	MessageType string `gorm:"size:50"`
	Status      string `gorm:"size:20"`
	Response    string
	SentAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
