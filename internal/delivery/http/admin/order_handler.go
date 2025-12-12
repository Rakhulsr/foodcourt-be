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
		"Orders":       resp.Data,
		"Total":        resp.Total,
		"Page":         resp.Page,
		"ActiveMenu":   "order",
		"Title":        "Daftar Pesanan",
		"csrf_token":   c.GetString("csrf_token"),
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

	updatedOrder, _ := h.orderUsecase.GetOrderByCode(code)

	data := gin.H{
		"Order":         updatedOrder,
		"PaymentStatus": updatedOrder.PaymentStatus,
		"CsrfToken":     c.GetString("csrf_token"),
	}

	c.HTML(http.StatusOK, "admin_order_row.html", data)
}

func (h *OrderHandler) SendNotification(c *gin.Context) {
	code := c.Param("code")

	err := h.orderUsecase.SendOrderNotificationToSeller(code)
	if err != nil {

		c.Header("HX-Trigger", `{"showMessage": {"type": "error", "message": "Gagal kirim WA: `+err.Error()+`"}}`)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "button_notify_sent.html", nil)
}
