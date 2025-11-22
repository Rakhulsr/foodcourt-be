package dto

import "github.com/Rakhulsr/foodcourt/internal/model"

type CreateOrderRequest struct {
	CustomerName  string            `json:"customer_name"`
	TableNumber   string            `json:"table_number"`
	PaymentMethod string            `json:"payment_method"`
	Items         []model.OrderItem `json:"items"`
}

type CreateOrderResponse struct {
	OrderCode  string `json:"order_code"`
	PaymentURL string `json:"payment_url" binding:"omitempty"`
	Message    string `json:"message"`
}

type CreateOrderItemRequest struct {
	MenuID   uint   `json:"menu_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Notes    string `json:"notes"`
}

type OrderListResponse struct {
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Data  []model.Order `json:"data"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type XenditCallbackRequest struct {
	ID          string  `json:"id"`
	ExternalID  string  `json:"external_id"`
	Status      string  `json:"status"`
	Amount      float64 `json:"amount"`
	PaidAmount  float64 `json:"paid_amount"`
	PayerEmail  string  `json:"payer_email"`
	Description string  `json:"description"`
	Updated     string  `json:"updated"`
	Created     string  `json:"created"`
}
