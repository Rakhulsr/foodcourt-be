package model

type OrderItem struct {
	ID              uint `gorm:"primaryKey"`
	OrderID         uint `gorm:"not null"`
	MenuID          uint `gorm:"not null"`
	BoothID         uint `gorm:"not null"`
	Quantity        int  `gorm:"not null"`
	PriceAtPurchase int  `gorm:"not null"`

	Menu Menu `gorm:"foreignKey:MenuID"`

	Booth Booth `gorm:"foreignKey:BoothID"`
	Notes string
}
