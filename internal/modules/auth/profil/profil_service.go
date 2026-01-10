package profil

import "starter-wahcah-be/internal/modules/auth/login"

type Service interface {
	GetProfile(userID uint) (*ProfileResponse, error)
}

type service struct {
	loginRepo login.Repository
}

func NewProfilService(loginRepo login.Repository) Service {
	return &service{
		loginRepo: loginRepo,
	}
}

func (s *service) GetProfile(userID uint) (*ProfileResponse, error) {

	data, err := s.loginRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &ProfileResponse{
		FullName: data.FirstName + " " + data.LastName,
		Email:    data.Email,
	}, nil
}
