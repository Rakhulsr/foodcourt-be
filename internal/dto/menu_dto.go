package dto

type MenuResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	IsAvailable bool   `json:"is_available"`
	Booth       struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"booth"`
}

type MenuListResponse struct {
	Menus []MenuResponse `json:"menus"`
	Total int            `json:"total"`
}
