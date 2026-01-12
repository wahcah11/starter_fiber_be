package login

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
