package auth

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	NIK      string `json:"nik" binding:"required"`
	RTId     string `json:"rt_id" binding:"required"`
}

type RegisterResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type LoginRequest struct {
	EmailOrNIK string `json:"email_or_nik" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Role    string `json:"role"`
	Message string `json:"message"`
}
