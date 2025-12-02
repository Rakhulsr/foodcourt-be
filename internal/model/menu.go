package model

type Menu struct {
	ID          uint   `gorm:"primaryKey"`
	BoothID     uint   `gorm:"not null"`
	Name        string `gorm:"size:100;not null"`
	Price       int    `gorm:"not null"`
	IsAvailable bool
	Booth       Booth  `gorm:"foreignKey:BoothID"`
	Category    string `gorm:"size:20;default:'makanan'"`
	Description string `gorm:"type:text"`
	ImagePath   string `gorm:"size:255"`
}
