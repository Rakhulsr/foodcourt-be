package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logUC usecase.LogUseCase
}

func NewLogHandler(uc usecase.LogUseCase) *LogHandler {
	return &LogHandler{logUC: uc}
}

func (h *LogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	logs, total, err := h.logUC.GetLogs(page, 20)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_log_list.html", gin.H{
		"Title":      "WhatsApp Logs",
		"ActiveMenu": "log",
		"Logs":       logs,
		"Total":      total,
		"Page":       page,
	})
}

func (h *LogHandler) TrackAndRedirect(c *gin.Context) {
	orderID, _ := strconv.ParseUint(c.Query("order_id"), 10, 32)
	boothID, _ := strconv.ParseUint(c.Query("booth_id"), 10, 32)
	phone := c.Query("phone")
	text := c.Query("text")

	go func() {
		h.logUC.RecordLog(uint(orderID), uint(boothID), phone)
	}()

	waURL := fmt.Sprintf("https://wa.me/%s?text=%s", phone, text)
	c.Redirect(http.StatusTemporaryRedirect, waURL)
}
