package dto

import "github.com/Rakhulsr/foodcourt/internal/model"

type CartItemCookie struct {
	MenuID   uint   `json:"menu_id"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
}

type AddToCartRequest struct {
	MenuID   uint `json:"menu_id" form:"menu_id" binding:"required"`
	Quantity int  `json:"quantity" form:"quantity"`
}

type CartView struct {
	Items       []CartItemView
	TotalAmount int
	TotalQty    int
}

type CartItemView struct {
	Menu     *model.Menu
	Quantity int
	SubTotal int
	Notes    string
}
