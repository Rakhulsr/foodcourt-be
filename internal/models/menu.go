package models

type Menu struct {
	ID          uint   `gorm:"primaryKey"`
	BoothID     uint   `gorm:"not null"`
	Name        string `gorm:"size:100;not null"`
	Price       int    `gorm:"not null"`
	IsAvailable bool   `gorm:"default:true"`
}
