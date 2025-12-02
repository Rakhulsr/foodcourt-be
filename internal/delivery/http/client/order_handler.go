package client

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/Rakhulsr/foodcourt/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(ou usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: ou}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBind(&req); err != nil {
		utils.SetFlash(c, "error", "Data pesanan tidak lengkap: "+err.Error())
		c.Redirect(http.StatusFound, "/checkout")
		return
	}

	cookie, _ := c.Cookie("user_cart")
	if cookie == "" {
		utils.SetFlash(c, "error", "Keranjang kosong")
		c.Redirect(http.StatusFound, "")
		return
	}

	jsonStr, _ := url.QueryUnescape(cookie)
	var cartItems []dto.CartItemCookie
	json.Unmarshal([]byte(jsonStr), &cartItems)

	if len(cartItems) == 0 {
		utils.SetFlash(c, "error", "Keranjang kosong")
		c.Redirect(http.StatusFound, "")
		return
	}

	for _, item := range cartItems {
		req.Items = append(req.Items, dto.CreateOrderItemRequest{
			MenuID:   item.MenuID,
			Quantity: item.Quantity,
			Notes:    item.Notes,
		})
	}

	res, err := h.orderUsecase.CreateOrder(req)
	if err != nil {

		utils.SetFlash(c, "error", "Gagal membuat pesanan: "+err.Error())
		c.Redirect(http.StatusFound, "/checkout")
		return
	}

	c.SetCookie("user_cart", "", -1, "/", "", false, false)
	c.SetCookie("temp_customer_name", "", -1, "/", "", false, false)
	c.SetCookie("temp_table_number", "", -1, "/", "", false, false)

	if res.PaymentURL != "" {
		c.Redirect(http.StatusFound, res.PaymentURL)
	} else {

		c.Redirect(http.StatusFound, "/order/success/"+res.OrderCode)
	}
}

func (h *OrderHandler) ShowSuccessPage(c *gin.Context) {
	orderCode := c.Param("code")

	order, err := h.orderUsecase.GetOrderByCode(orderCode)
	if err != nil {
		c.String(http.StatusNotFound, "Order tidak ditemukan")
		return
	}

	c.HTML(http.StatusOK, "order_success.html", gin.H{
		"Title":        "Pesanan Berhasil",
		"Order":        order,
		"FlashMessage": c.GetString("FlashMessage"),
		"FlashType":    c.GetString("FlashType"),
	})
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
