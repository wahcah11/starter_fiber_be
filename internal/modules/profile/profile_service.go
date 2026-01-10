package profile

type Service interface {
	GetByID(id uint) (User, error)
}

type service struct {
	repo Repository
}

func NewProfileService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetByID(id uint) (User, error) {
	return s.repo.GetByID(id)
}
