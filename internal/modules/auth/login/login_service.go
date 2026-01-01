package login

import (
	"errors"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Authenticate(req LoginRequest) (*LoginResponse, error)
	RegisterUser(email, password, firstname, lastname string) error // Helper buat bikin user
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

	token, _ := util.GenerateToken(user.ID)

	return &LoginResponse{
		Token:     token,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}, nil
}

func (s *service) RegisterUser(email, password, firstname, lastname string) error {
	hashed, _ := util.HashPassword(password)

	user := User{
		Email:     email,
		Password:  hashed,
		Firstname: firstname,
		Lastname:  lastname,
	}

	return s.repo.CreateUser(&user)
}
