package dto

import "time"

type BoothCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	WhatsApp string `json:"whatsapp" binding:"required"`
	IsActive bool   `json:"is_active" binding:"omitempty"`
}

type BoothUpdateRequest struct {
	Name     string `json:"name" binding:"omitempty"`
	WhatsApp string `json:"whatsapp" binding:"omitempty"`
	IsActive bool   `json:"is_active" binding:"omitempty"`
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
