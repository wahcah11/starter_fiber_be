package login

type LoginRequest struct {
    Firstname string `json:"firstname" validate:"required"`
    Lastname  string `json:"lastname" validate:"required"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
    Token     string `json:"token"`
    Firstname string `json:"firstname"`
    Lastname  string `json:"lastname"`
}

