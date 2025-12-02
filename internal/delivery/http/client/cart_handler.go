package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/Rakhulsr/foodcourt/utils"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	menuUC usecase.MenuUseCase
}

func NewCartHandler(muc usecase.MenuUseCase) *CartHandler {
	return &CartHandler{menuUC: muc}
}

const (
	CartCookieName     = "user_cart"
	CustomerNameCookie = "temp_customer_name"
	TableNumberCookie  = "temp_table_number"
)

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

	for i, item := range cartItems {
		if item.MenuID == req.MenuID {
			cartItems[i].Quantity += req.Quantity
			found = true
			break
		}
	}
	if !found {
		cartItems = append(cartItems, dto.CartItemCookie{
			MenuID: req.MenuID, Quantity: req.Quantity,
		})
	}

	h.saveCartToCookie(c, cartItems)

	totalQty := 0
	var itemNames []string

	for _, item := range cartItems {
		totalQty += item.Quantity

		menu, err := h.menuUC.GetByID(item.MenuID)
		if err == nil {

			itemNames = append(itemNames, fmt.Sprintf("%s x%d", menu.Name, item.Quantity))
		}
	}

	summaryText := ""
	if len(itemNames) > 0 {
		summaryText = strings.Join(itemNames, ", ")

		if len(summaryText) > 40 {
			summaryText = summaryText[:37] + "..."
		}
	}

	c.HTML(http.StatusOK, "cart_add_response.html", gin.H{
		"TotalQty":    totalQty,
		"CartSummary": summaryText,
	})
}

func (h *CartHandler) ShowCart(c *gin.Context) {
	cookieItems := h.getCartFromCookie(c)

	customerName, _ := c.Cookie(CustomerNameCookie)
	tableNumber, _ := c.Cookie(TableNumberCookie)

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
			"Notes":     item.Notes,
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
		"CustomerName": customerName,
		"TableNumber":  tableNumber,
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
	c.SetCookie(CartCookieName, cookieValue, 3600, "/", "", false, false)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	menuID, _ := strconv.Atoi(c.PostForm("menu_id"))
	action := c.PostForm("action")
	note := c.PostForm("note")

	cartItems := h.getCartFromCookie(c)
	var newItems []dto.CartItemCookie

	var currentItem dto.CartItemCookie
	var itemSubTotal int

	for _, item := range cartItems {
		if item.MenuID == uint(menuID) {

			if action == "update_note" {
				item.Notes = note
				newItems = append(newItems, item)
			} else if action == "increase" {
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

func (h *CartHandler) ProceedCheckout(c *gin.Context) {
	customerName := c.PostForm("customer_name")
	tableNumber := c.PostForm("table_number")

	if customerName == "" {
		utils.SetFlash(c, "error", "Nama pemesan wajib diisi!")
		c.Redirect(http.StatusFound, "/cart")
		return
	}

	c.SetCookie("temp_customer_name", customerName, 3600, "/", "", false, false)
	c.SetCookie(TableNumberCookie, tableNumber, 3600, "/", "", false, false)

	c.Redirect(http.StatusFound, "/checkout")
}

func (h *CartHandler) ShowCheckoutPage(c *gin.Context) {
	cookieItems := h.getCartFromCookie(c)

	customerName, _ := c.Cookie(CustomerNameCookie)
	tableNumber, _ := c.Cookie(TableNumberCookie)

	if len(cookieItems) == 0 {
		c.Redirect(http.StatusFound, "")
		return
	}

	if customerName == "" {
		utils.SetFlash(c, "error", "Sesi habis, silakan isi nama kembali")
		c.Redirect(http.StatusFound, "/cart")
		return
	}

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

			"Name":     menu.Name,
			"Price":    menu.Price,
			"Quantity": item.Quantity,
			"SubTotal": subTotal,
			"Notes":    item.Notes,
		})
	}

	c.HTML(http.StatusOK, "client_checkout.html", gin.H{
		"Title":        "Konfirmasi Pesanan",
		"CartItems":    finalItems,
		"TotalAmount":  totalAmount,
		"TotalQty":     totalQty,
		"CustomerName": customerName,
		"TableNumber":  tableNumber,
		"csrf_token":   c.GetString("csrf_token"),
	})
}
