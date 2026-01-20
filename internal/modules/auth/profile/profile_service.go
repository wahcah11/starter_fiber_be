package profile

import (
	"errors"
	"strings"

	"starter-wahcah-be/internal/modules/auth/login"
)

type Service interface {
	GetProfile(userID uint) (*ProfileResponse, error)
}

type service struct {
	loginRepo login.Repository
}

func NewProfileService(loginRepo login.Repository) Service {
	return &service{loginRepo: loginRepo}
}

func (s *service) GetProfile(userID uint) (*ProfileResponse, error) {
	user, err := s.loginRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// pakai kolom Nama kalau ada (lebih cocok dengan aturan kamu sebelumnya)
	fullName := strings.TrimSpace(user.Nama)
	if fullName == "" {
		fullName = strings.TrimSpace(user.FirstName + " " + user.LastName)
	}

	return &ProfileResponse{
		Nama:  fullName,
		Email: user.Email,
	}, nil
}
