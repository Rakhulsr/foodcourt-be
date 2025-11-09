package model

import "time"

type Booth struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	WhatsApp  string `gorm:"size:20;not null"`
	IsActive  bool   `gorm:"default:true"`
	Menus     []Menu `gorm:"foreignKey:BoothID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
