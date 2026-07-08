package dto

type RegisterRequest struct {
	Name        string `json:"name"         binding:"required,min=2"`
	Email       string `json:"email"        binding:"required,email"`
	Password    string `json:"password"     binding:"required,min=6"`
	Phone       string `json:"phone"`
	OperatorKey string `json:"operator_key"` // opcional — si coincide → operador
}
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
