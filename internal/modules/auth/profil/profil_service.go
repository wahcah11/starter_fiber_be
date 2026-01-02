package profil

import (
    "errors"
    "starter-wahcah-be/internal/modules/auth/login"
)

type Service interface {
    GetProfile(userID uint) (*ProfilResponse, error)
}

type service struct {
    loginRepo login.Repository
}

func NewProfilService(loginRepo login.Repository) Service {
    return &service{loginRepo}
}

func (s *service) GetProfile(userID uint) (*ProfilResponse, error) {
    user, err := s.loginRepo.FindByID(userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    return &ProfilResponse{
        Fullname: user.Firstname + " " + user.Lastname,
        Email:    user.Email,
    }, nil
}
