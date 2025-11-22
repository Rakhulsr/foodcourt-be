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

type MenuHandler struct {
	menuUc  usecase.MenuUseCase
	boothUC usecase.BoothUseCase
}

func NewMenuHandler(uc usecase.MenuUseCase, bu usecase.BoothUseCase) *MenuHandler {
	return &MenuHandler{menuUc: uc, boothUC: bu}
}

func (h *MenuHandler) ListActive(c *gin.Context) {
	category := c.Query("category")
	keyword := c.Query("keyword")

	var resp *dto.MenuListResponse
	var err error

	if keyword != "" {
		resp, err = h.menuUc.FindByName(keyword)
	} else if category != "" {
		resp, err = h.menuUc.FindByCategory(category)
	} else {
		resp, err = h.menuUc.ListActive()
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "menu_list_client.html", gin.H{
		"Menus":        resp.Menus,
		"Total":        resp.Total,
		"Category":     category,
		"Keyword":      keyword,
		"FlashMessage": c.GetString("FlashMessage"),
		"FlashType":    c.GetString("FlashType"),
	})
}

func (h *MenuHandler) ListByBooth(c *gin.Context) {
	boothIDStr := c.Param("booth_id")
	boothID, _ := strconv.ParseUint(boothIDStr, 10, 32)

	resp, err := h.menuUc.ListActiveByBoothID(uint(boothID))
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Booth not found or empty"})
		return
	}

	c.HTML(http.StatusOK, "menu_by_booth.html", gin.H{
		"Menus":   resp.Menus,
		"BoothID": boothID,
	})
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	resp, err := h.menuUc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")

	resp, err := h.menuUc.FindByName(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) FilterByCategory(c *gin.Context) {
	category := c.Query("category")

	resp, err := h.menuUc.FindByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) ClientHome(c *gin.Context) {
	keyword := c.Query("keyword")

	boothsResp, err := h.boothUC.ListActive()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var menusResp *dto.MenuListResponse

	if keyword != "" {
		menusResp, err = h.menuUc.FindByName(keyword)
	} else {
		menusResp, err = h.menuUc.ListActive()
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	cookie, _ := c.Cookie("user_cart")
	totalQty := 0
	if cookie != "" {
		jsonStr, _ := url.QueryUnescape(cookie)
		var items []dto.CartItemCookie
		json.Unmarshal([]byte(jsonStr), &items)
		for _, item := range items {
			totalQty += item.Quantity
		}
	}

	c.HTML(http.StatusOK, "client_home.html", gin.H{
		"Title":        "Beranda",
		"Booths":       boothsResp.Booths,
		"Menus":        menusResp.Menus,
		"Keyword":      keyword,
		"FlashMessage": c.GetString("FlashMessage"),
		"FlashType":    c.GetString("FlashType"),
		"TotalQty":     totalQty,
		"csrf_token":   c.GetString("csrf_token"),
		"ActiveTab":    "home",
	})
}
