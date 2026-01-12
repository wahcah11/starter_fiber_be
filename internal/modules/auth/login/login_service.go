package login

import (
	"errors"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (User, error)
	RegisterUser(email, password string) error
	GetByID(id uint) (User, error)
}

type service struct {
	repo Repository
}

func NewLoginService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(req RegisterRequest) error {
	hashed, _ := util.HashPassword(req.Password)

	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashed,
	}

	return s.repo.CreateUser(&user)
}

func (s *service) Login(req LoginRequest) (User, error) {
	user, err := s.repo.Login(req)
	if err != nil {
		return user, errors.New("invalid email or password")
	}

	// cek password hash
	if !util.CheckPasswordHash(req.Password, user.Password) {
		return user, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *service) GetByID(id uint) (User, error) {
	return s.repo.GetByID(id)
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

func (s *service) RegisterUser(email, password string) error {
	hashed, _ := util.HashPassword(password)
	user := User{Email: email, Password: hashed}
	return s.repo.CreateUser(&user)
}
