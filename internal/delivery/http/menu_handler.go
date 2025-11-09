package http

import (
	"net/http"
	"strconv"

	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	uc usecase.MenuUseCase
}

func NewMenuHandler(uc usecase.MenuUseCase) *MenuHandler {
	return &MenuHandler{uc: uc}
}

func (h *MenuHandler) List(c *gin.Context) {
	resp, err := h.uc.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	resp, err := h.uc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) ListByBooth(c *gin.Context) {
	boothIDStr := c.Param("booth_id")
	boothID, err := strconv.ParseUint(boothIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Booth ID"})
		return
	}

	resp, err := h.uc.ListActiveByBoothID(uint(boothID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")

	resp, err := h.uc.FindByName(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) FilterByCategory(c *gin.Context) {
	category := c.Query("category")

	resp, err := h.uc.FindByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
