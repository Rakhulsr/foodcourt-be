package model

import "time"

type WhatsAppLog struct {
	ID          uint      `gorm:"primaryKey"`
	OrderID     *uint     `gorm:"index"`
	Order       *Order    `gorm:"foreignKey:OrderID"`
	BoothID     *uint     `gorm:"index"`
	Booth       *Booth    `gorm:"foreignKey:BoothID"`
	MessageType string    `gorm:"size:50"`
	Status      string    `gorm:"size:20"`
	Response    string    `gorm:"type:text"`
	SentAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
