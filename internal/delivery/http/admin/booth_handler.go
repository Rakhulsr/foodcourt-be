package admin

import (
	"net/http"
	"strconv"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/Rakhulsr/foodcourt/utils"
	"github.com/gin-gonic/gin"
)

type BoothHandler struct {
	uc usecase.BoothUseCase
}

func NewBoothHandler(uc usecase.BoothUseCase) *BoothHandler {
	return &BoothHandler{uc: uc}
}

func (h *BoothHandler) AdminList(c *gin.Context) {
	resp, err := h.uc.ListAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_booth_list.html", gin.H{
		"Booths":       resp.Booths,
		"Title":        "Kelola Booth",
		"ActiveMenu":   "booth",
		"FlashMessage": c.GetString("FlashMessage"),
		"FlashType":    c.GetString("FlashType"),
		"csrf_token":   c.GetString("csrf_token"),
	})
}

func (h *BoothHandler) ShowCreateForm(c *gin.Context) {

	c.HTML(http.StatusOK, "admin_booth_form.html", gin.H{
		"Type":       "create",
		"Title":      "Tambah Booth Baru",
		"ActiveMenu": "booth",
		"csrf_token": c.GetString("csrf_token"),
	})
}

func (h *BoothHandler) Create(c *gin.Context) {
	var req dto.BoothCreateRequest

	if err := c.ShouldBind(&req); err != nil {
		println("❌ VALIDASI GAGAL:", err.Error())

		c.HTML(http.StatusBadRequest, "admin_booth_form.html", gin.H{
			"Error": err.Error(),
			"Type":  "create",
			"Title": "Tambah Booth Baru",
		})
		return
	}

	if c.PostForm("is_active") == "on" {
		req.IsActive = true
	} else {
		req.IsActive = false
	}

	_, err := h.uc.Create(req)
	if err != nil {
		println("❌ ERROR CREATE BOOTH:", err.Error())

		c.HTML(http.StatusInternalServerError, "admin_booth_form.html", gin.H{
			"Error": err.Error(),
			"Type":  "create",
			"Title": "Tambah Booth Baru",
			"Data":  req,
		})
		return
	}

	utils.SetFlash(c, "success", "Booth berhasil dibuat!")
	c.Header("HX-Redirect", "/api/admin/booths")
	c.Status(http.StatusCreated)
}

func (h *BoothHandler) ShowEditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	booth, err := h.uc.GetByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/api/admin/booths")
		return
	}

	c.HTML(http.StatusOK, "admin_booth_form.html", gin.H{
		"Type":       "edit",
		"Title":      "Edit Booth: " + booth.Name,
		"ActiveMenu": "booth",
		"Data":       booth,
		"csrf_token": c.GetString("csrf_token"),
	})
}

func (h *BoothHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req dto.BoothUpdateRequest
	if err := c.ShouldBind(&req); err != nil {

		c.Redirect(http.StatusFound, "/api/admin/booths")
		return
	}

	if c.PostForm("is_active") == "on" {
		req.IsActive = true
	} else {
		req.IsActive = false
	}

	updatedBooth, err := h.uc.Update(uint(id), req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "booth_row.html", updatedBooth)
}

func (h *BoothHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.uc.Delete(uint(id)); err != nil {

		c.String(http.StatusInternalServerError, "Gagal menghapus")
		return
	}

	c.HTML(http.StatusOK, "flash.html", gin.H{"Type": "success", "Message": "Booth dihapus!"})
}
