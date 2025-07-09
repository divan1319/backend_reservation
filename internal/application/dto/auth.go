package dto

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}
