package profile

// import "starter-wahcah-be/internal/modules/auth/login"

type Service interface {
	GetByID(id uint) (User, error)
}

type service struct {
	repo Repository
}

func NewProfileService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(id uint) (User, error) {
	return s.repo.GetByID(id)
}
