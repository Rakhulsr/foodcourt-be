package dto

import "time"

type BoothCreateRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	WhatsApp string `json:"whatsapp" form:"whatsapp" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type BoothUpdateRequest struct {
	Name     string `json:"name" form:"name"`
	WhatsApp string `json:"whatsapp" form:"whatsapp"`
	IsActive bool   `json:"is_active"`
}

type BoothResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	WhatsApp  string    `json:"whatsapp"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BoothListResponse struct {
	Booths []BoothResponse `json:"booths"`
	Total  int             `json:"total"`
}
