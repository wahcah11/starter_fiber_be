package profile

import "starter-wahcah-be/internal/modules/auth/login"

type Service interface {
	GetProfile(userID uint) (*ProfileResponse, error)
}

type service struct {
	repo login.Repository
}

func NewProfileService(repo login.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile(userID uint) (*ProfileResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	fullName := user.FirstName + " " + user.LastName
	return &ProfileResponse{
		Name:  fullName,
		Email: user.Email,
	}, nil
}
