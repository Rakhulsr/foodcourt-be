package dto

type MenuCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    int    `json:"price" binding:"required,min=1000"`
	Category string `json:"category" binding:"required,oneof=makanan minuman"`
	BoothID  uint   `json:"booth_id" binding:"required"`
}

type MenuUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Price       int    `json:"price" binding:"omitempty,min=1000"`
	Category    string `json:"category" binding:"omitempty,oneof=makanan minuman"`
	IsAvailable bool   `json:"is_available" binding:"omitempty"`
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
	Menus []MenuResponse `json:"menus"`
	Total int            `json:"total"`
}
