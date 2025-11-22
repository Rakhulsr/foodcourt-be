package dto

type MenuCreateRequest struct {
	BoothID     uint   `json:"booth_id" form:"booth_id" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Price       int    `json:"price" form:"price" binding:"required"`
	Category    string `json:"category" form:"category"`
	IsAvailable bool   `json:"is_available"`
}

type MenuUpdateRequest struct {
	Name        string `json:"name" form:"name"`
	Price       int    `json:"price" form:"price"`
	Category    string `json:"category" form:"category"`
	IsAvailable bool   `json:"is_available"`
}

type MenuResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	IsAvailable bool   `json:"is_available"`
	Category    string `json:"category"`
	ImagePath   string `json:"image_path"`
	Booth       struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"booth"`
}

type MenuListResponse struct {
	Total int            `json:"total"`
	Menus []MenuResponse `json:"menus"`
}
