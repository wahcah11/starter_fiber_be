package profil

import "errors"

type Service interface {
	GetProfile(userID uint) (*ProfilResponse, error)
}

type service struct {
	repo Repository
}

func NewProfilService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetProfile(userID uint) (*ProfilResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &ProfilResponse{
		Fullname: user.Firstname + " " + user.Lastname,
		Email:    user.Email,
	}, nil
}
