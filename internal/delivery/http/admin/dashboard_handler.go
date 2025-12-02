package admin

import (
	"net/http"

	"github.com/Rakhulsr/foodcourt/internal/repository"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	orderRepo repository.OrderRepository
}

func NewDashboardHandler(or repository.OrderRepository) *DashboardHandler {
	return &DashboardHandler{orderRepo: or}
}

func (h *DashboardHandler) Dashboard(c *gin.Context) {

	income, err := h.orderRepo.GetTotalIncomeToday()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	totalOrder, err := h.orderRepo.CountOrdersToday()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	orders, err := h.orderRepo.FindOrdersToday()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"Title":      "Dashboard",
		"ActiveMenu": "dashboard",
		"AdminName":  "Admin",
		"Income":     income,
		"Outcome":    0,
		"TotalOrder": totalOrder,
		"Orders":     orders,

		"csrf_token": c.GetString("csrf_token"),
	})
}
