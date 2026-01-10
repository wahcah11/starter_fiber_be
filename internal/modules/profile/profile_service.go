package profile

import "starter-wahcah-be/internal/modules/auth/login"

type Service interface {
	GetByID(id uint) (login.User, error)
}

type service struct {
	loginRepo login.Repository
}

func NewProfileService(repo login.Repository) Service {
	return &service{loginRepo: repo}
}

func (s *service) GetByID(id uint) (login.User, error) {
	return s.loginRepo.GetByID(id)
}
