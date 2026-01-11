package login

import (
	"errors"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Authenticate(req LoginRequest) (*LoginResponse, error)
	RegisterUser(req RegisterRequest) error // Helper buat bikin user
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

	token, err := util.GenerateToken(user.ID)

	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	return &LoginResponse{
		Token:     token,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (s *service) RegisterUser(req RegisterRequest) error {
	hashed, _ := util.HashPassword(req.Password)
	
	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashed,
	}
	return s.repo.CreateUser(&user)
}