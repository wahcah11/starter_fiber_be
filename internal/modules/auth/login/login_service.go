package login

import (
	"errors"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Authenticate(req LoginRequest) (*LoginResponse, error)
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
	return &LoginResponse{Token: token}, nil
}
