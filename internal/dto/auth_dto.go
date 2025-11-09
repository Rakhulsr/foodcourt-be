package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Admin struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		FullName string `json:"full_name"`
	} `json:"admin"`
	Message string `json:"message"`
}
