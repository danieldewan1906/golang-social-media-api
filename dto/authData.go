package dto

type RegisterRequest struct {
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName"`
	Password   string `json:"password" validate:"required"`
	Repassword string `json:"repassword" validate:"required"`
	Role       string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
