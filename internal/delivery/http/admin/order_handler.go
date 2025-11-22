package admin

import (
	"net/http"
	"strconv"

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

func (h *OrderHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	status := c.Query("status")

	resp, err := h.orderUsecase.ListOrders(page, limit, status)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_order_list.html", gin.H{
		"Orders":     resp.Data,
		"Total":      resp.Total,
		"Page":       resp.Page,
		"ActiveMenu": "order",
		"Title":      "Daftar Pesanan",

		"ActiveStatus": status,
	})
}

func (h *OrderHandler) AdminUpdateStatus(c *gin.Context) {
	code := c.Param("code")
	var req dto.UpdateStatusRequest

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid Data")
		return
	}

	err := h.orderUsecase.UpdateOrderStatus(code, req.Status)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data := gin.H{
		"PaymentStatus": "paid",
		"OrderStatus":   req.Status,
	}

	c.HTML(http.StatusOK, "status_badge.html", data)
}
