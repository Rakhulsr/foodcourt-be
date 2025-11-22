package client

import (
	"net/http"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(ou usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: ou}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBind(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.orderUsecase.CreateOrder(req)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	if res.PaymentURL != "" {
		c.Redirect(http.StatusFound, res.PaymentURL)
	} else {

		c.Redirect(http.StatusFound, "/order/success/"+res.OrderCode)
	}
}

func (h *OrderHandler) ShowOrderPage(c *gin.Context) {
	orderCode := c.Param("code")

	if orderCode == "" {
		c.String(http.StatusBadRequest, "missing order code")
		return
	}

	orderData, err := h.orderUsecase.GetOrderByCode(orderCode)
	if err != nil {
		c.String(http.StatusNotFound, "order not found")
		return
	}

	c.HTML(http.StatusOK, "order_success.html", orderData)

}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	code := c.Param("code")
	order, err := h.orderUsecase.GetOrderByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) HandleXenditWebhook(c *gin.Context) {

	var req dto.XenditCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := h.orderUsecase.ProcessXenditCallback(req)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
