package client

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	menuUC usecase.MenuUseCase
}

func NewCartHandler(muc usecase.MenuUseCase) *CartHandler {
	return &CartHandler{menuUC: muc}
}

const CartCookieName = "user_cart"

func (h *CartHandler) AddToCart(c *gin.Context) {
	var req dto.AddToCartRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	cartItems := h.getCartFromCookie(c)
	found := false

	totalQty := 0

	for i, item := range cartItems {
		if item.MenuID == req.MenuID {
			cartItems[i].Quantity += req.Quantity
			found = true
		}

	}
	if !found {
		cartItems = append(cartItems, dto.CartItemCookie{
			MenuID: req.MenuID, Quantity: req.Quantity,
		})
	}

	totalQty = 0
	for _, item := range cartItems {
		totalQty += item.Quantity
	}

	h.saveCartToCookie(c, cartItems)

	c.HTML(http.StatusOK, "cart_add_response.html", gin.H{
		"TotalQty": totalQty,
	})
}

func (h *CartHandler) ShowCart(c *gin.Context) {
	cookieItems := h.getCartFromCookie(c)

	var finalItems []map[string]interface{}
	totalAmount := 0
	totalQty := 0

	for _, item := range cookieItems {
		menu, err := h.menuUC.GetByID(item.MenuID)
		if err != nil {
			continue
		}

		subTotal := menu.Price * item.Quantity
		totalAmount += subTotal
		totalQty += item.Quantity

		finalItems = append(finalItems, map[string]interface{}{
			"MenuID":    menu.ID,
			"Name":      menu.Name,
			"Price":     menu.Price,
			"ImagePath": menu.ImagePath,
			"BoothName": menu.Booth.Name,
			"Quantity":  item.Quantity,
			"SubTotal":  subTotal,
		})
	}

	c.HTML(http.StatusOK, "client_cart.html", gin.H{
		"Title":        "Keranjang Pesanan",
		"CartItems":    finalItems,
		"TotalAmount":  totalAmount,
		"TotalQty":     totalQty,
		"ActiveTab":    "cart",
		"csrf_token":   c.GetString("csrf_token"),
		"FlashMessage": c.GetString("FlashMessage"),
		"FlashType":    c.GetString("FlashType"),
	})
}

func (h *CartHandler) getCartFromCookie(c *gin.Context) []dto.CartItemCookie {
	cookie, err := c.Cookie(CartCookieName)
	if err != nil || cookie == "" {
		return []dto.CartItemCookie{}
	}
	jsonStr, _ := url.QueryUnescape(cookie)
	var items []dto.CartItemCookie
	json.Unmarshal([]byte(jsonStr), &items)
	return items
}

func (h *CartHandler) saveCartToCookie(c *gin.Context, items []dto.CartItemCookie) {
	jsonData, _ := json.Marshal(items)
	cookieValue := url.QueryEscape(string(jsonData))
	c.SetCookie(CartCookieName, cookieValue, 3600*24*7, "/", "", false, false)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	menuID, _ := strconv.Atoi(c.PostForm("menu_id"))
	action := c.PostForm("action")

	cartItems := h.getCartFromCookie(c)
	var newItems []dto.CartItemCookie

	var currentItem dto.CartItemCookie
	var itemSubTotal int

	for _, item := range cartItems {
		if item.MenuID == uint(menuID) {
			if action == "increase" {
				item.Quantity++
				newItems = append(newItems, item)
			} else if action == "decrease" {
				if item.Quantity > 1 {
					item.Quantity--
					newItems = append(newItems, item)
				} else {

					item.Quantity = 1
					newItems = append(newItems, item)
				}
			} else if action == "delete" {

				currentItem = item
				continue
			}
			currentItem = item
		} else {
			newItems = append(newItems, item)
		}
	}

	h.saveCartToCookie(c, newItems)

	totalAmount := 0
	totalQty := 0

	for _, item := range newItems {
		menu, _ := h.menuUC.GetByID(item.MenuID)
		subTotal := menu.Price * item.Quantity
		totalAmount += subTotal
		totalQty += item.Quantity

		if item.MenuID == uint(menuID) {
			itemSubTotal = subTotal

			if action != "delete" {
				currentItem.Quantity = item.Quantity
			}
		}
	}

	c.HTML(http.StatusOK, "cart_update.html", gin.H{
		"Item":         currentItem,
		"ItemSubTotal": itemSubTotal,
		"TotalAmount":  totalAmount,
		"TotalQty":     totalQty,
		"Action":       action,
		"CsrfToken":    c.PostForm("csrf_token"),
	})
}
