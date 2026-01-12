package profil

import "starter-wahcah-be/internal/modules/auth/login"

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
		return nil, err
	}

	return &ProfilResponse{
		Fullname: user.Firstname + " " + user.Lastname,
		Email:    user.Email,
	}, nil
}
