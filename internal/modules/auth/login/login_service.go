package login

import (
	"errors"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Authenticate(req LoginRequest) (*LoginResponse, error)
	RegisterUser(email, password string) error // Helper buat bikin user
}

type service struct {
	repo Repository
}

func NewLoginService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Authenticate(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}
	type userValue struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
	}

	token, _ := util.GenerateToken(user.ID)
	return &LoginResponse{Token: token, User: userValue{FirstName: user.FirstName, LastName: user.LastName} }, nil
}

func (s *service) RegisterUser(email, password string) error {
	hashed, _ := util.HashPassword(password)
	user := User{Email: email, Password: hashed}
	return s.repo.CreateUser(&user)
}


