package client

import (
	"net/http"

	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/gin-gonic/gin"
)

type BoothHandler struct {
	uc usecase.BoothUseCase
}

func NewBoothHandler(uc usecase.BoothUseCase) *BoothHandler {
	return &BoothHandler{uc: uc}
}

func (h *BoothHandler) List(c *gin.Context) {
	resp, err := h.uc.ListActive()
	if err != nil {

		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "booth_list_client.html", gin.H{
		"Booths": resp.Booths,
		"Total":  resp.Total,
	})
}
