package admin

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Rakhulsr/foodcourt/internal/dto"
	"github.com/Rakhulsr/foodcourt/internal/usecase"
	"github.com/Rakhulsr/foodcourt/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MenuHandler struct {
	menuUC  usecase.MenuUseCase
	boothUC usecase.BoothUseCase
}

func NewMenuHandler(uc usecase.MenuUseCase, bu usecase.BoothUseCase) *MenuHandler {
	return &MenuHandler{menuUC: uc, boothUC: bu}
}

func (h *MenuHandler) ListAll(c *gin.Context) {
	resp, err := h.menuUC.ListAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_menu_list.html", gin.H{
		"Menus":      resp.Menus,
		"Title":      "Kelola Menu",
		"ActiveMenu": "menu",

		"FlashMessage": c.GetString("FlashMessage"),
		"FlashType":    c.GetString("FlashType"),
		"csrf_token":   c.GetString("csrf_token"),
	})
}

func (h *MenuHandler) ShowCreateForm(c *gin.Context) {

	booths, _ := h.boothUC.ListActive()

	c.HTML(http.StatusOK, "admin_menu_form.html", gin.H{
		"Type":       "create",
		"Title":      "Tambah Menu Baru",
		"Booths":     booths.Booths,
		"csrf_token": c.GetString("csrf_token"),
	})
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.MenuCreateRequest

	if err := c.ShouldBind(&req); err != nil {
		booths, _ := h.boothUC.ListActive()
		c.HTML(http.StatusBadRequest, "admin_menu_form.html", gin.H{
			"Error":  err.Error(),
			"Type":   "create",
			"Booths": booths.Booths,
			"Title":  "Tambah Menu Baru",
		})
		return
	}

	if c.PostForm("is_available") == "on" {
		req.IsAvailable = true
	} else {
		req.IsAvailable = false
	}

	imagePath := ""
	file, err := c.FormFile("image")

	if err == nil {

		ext := filepath.Ext(file.Filename)
		filename := uuid.New().String() + ext

		saveDir := filepath.Join("public", "uploads", "menu")
		savePath := filepath.Join(saveDir, filename)

		if _, err := os.Stat(saveDir); os.IsNotExist(err) {
			os.MkdirAll(saveDir, os.ModePerm)
		}

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Gagal save gambar: " + err.Error()})
			return
		}

		imagePath = "/uploads/menu/" + filename
	}

	_, err = h.menuUC.Create(req, imagePath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	utils.SetFlash(c, "success", "Menu berhasil ditambahkan!")

	c.Header("HX-Redirect", "/api/admin/menus")
	c.Status(http.StatusCreated)
}

func (h *MenuHandler) ShowEditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	menu, err := h.menuUC.GetByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/api/admin/menus")
		return
	}

	boothsResp, err := h.boothUC.ListActive()
	if err != nil {
		log.Printf("failed to get booths for menus: %v", err)
		c.Redirect(http.StatusNotFound, "/api/admin/menus")
		return
	}

	c.HTML(http.StatusOK, "admin_menu_form.html", gin.H{
		"Type":       "edit",
		"Title":      "Edit Menu: " + menu.Name,
		"Data":       menu,
		"Booths":     boothsResp.Booths,
		"csrf_token": c.GetString("csrf_token"),
	})
}

func (h *MenuHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req dto.MenuUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid Input: "+err.Error())
		c.Redirect(http.StatusFound, "/api/admin/menus")
		return
	}

	if c.PostForm("is_available") == "on" {
		req.IsAvailable = true
	} else {
		req.IsAvailable = false
	}

	imagePath := ""
	file, err := c.FormFile("image")

	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := uuid.New().String() + ext

		saveDir := filepath.Join("public", "uploads", "menu")
		savePath := filepath.Join(saveDir, filename)
		if _, err := os.Stat(saveDir); os.IsNotExist(err) {
			os.MkdirAll(saveDir, os.ModePerm)
		}

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.String(http.StatusInternalServerError, "Gagal upload gambar")
			return
		}

		imagePath = "/uploads/menu/" + filename
	}

	updatedMenu, err := h.menuUC.Update(uint(id), req, imagePath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "menu_row.html", updatedMenu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.menuUC.Delete(uint(id)); err != nil {

		c.HTML(http.StatusOK, "flash.html", gin.H{
			"Type":    "error",
			"Message": "Gagal menghapus: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "flash.html", gin.H{
		"Type":    "success",
		"Message": "Menu berhasil dihapus!",
	})
}
