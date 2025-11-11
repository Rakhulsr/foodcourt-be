package http

import (
	"net/http"

	"github.com/Rakhulsr/foodcourt/internal/model"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/Rakhulsr/foodcourt/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type OrderHandler struct {
	DB             *gorm.DB
	PaymentUsecase *usecase.PaymentUsecase
}

func NewOrderHandler(db *gorm.DB, ps *usecase.PaymentUsecase) *OrderHandler {
	return &OrderHandler{db, ps}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerName  string            `json:"customer_name"`
		TableNumber   string            `json:"table_number"`
		PaymentMethod string            `json:"payment_method"`
		Items         []model.OrderItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := model.Order{
		OrderCode:     "ORD-" + utils.RandomString(10),
		CustomerName:  req.CustomerName,
		TableNumber:   req.TableNumber,
		PaymentMethod: req.PaymentMethod,
		OrderStatus:   "pending",
		PaymentStatus: "pending",
	}

	total := 0
	for i := range req.Items {
		total += req.Items[i].PriceAtPurchase * req.Items[i].Quantity
	}
	order.TotalAmount = total
	order.Items = req.Items

	if err := h.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.PaymentMethod == "cash" {
		c.JSON(http.StatusOK, gin.H{
			"order_code": order.OrderCode,
			"message":    "Order created successfully. Please pay at cashier.",
		})
		return
	}

	inv, _ := h.PaymentUsecase.CreateInvoice(order)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	h.DB.Model(&order).Updates(map[string]interface{}{
		"XenditInvoiceID": inv.Id,
		"InvoiceURL":      inv.InvoiceUrl,
	})

	c.JSON(http.StatusOK, gin.H{
		"order_code":  order.OrderCode,
		"payment_url": inv.InvoiceUrl,
	})
}
