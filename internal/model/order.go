package model

import "time"

type Order struct {
	ID            uint   `gorm:"primaryKey"`
	OrderCode     string `gorm:"unique;size:20;not null"`
	CustomerName  string `gorm:"size:100;not null"`
	TableNumber   string `gorm:"size:10"`
	TotalAmount   int    `gorm:"not null"`
	PaymentMethod string `gorm:"type:enum('qris','cash');not null"`
	PaymentStatus string `gorm:"type:enum('pending','paid');default:'pending'"`
	OrderStatus   string `gorm:"type:enum('pending','confirmed','preparing','ready','completed');default:'pending'"`

	XenditInvoiceID string `gorm:"size:100"`
	InvoiceURL      string `gorm:"size:255"`

	AdminTransferred bool `gorm:"default:false"`
	AdminNote        string

	CreatedAt time.Time
	UpdatedAt time.Time

	Items []OrderItem `gorm:"foreignKey:OrderID"`
}
