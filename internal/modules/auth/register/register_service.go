package register

import (
	"errors"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Register(req RegisterRequest) (*RegisterResponse, error)
}

type service struct {
	repo Repository
}

func NewRegisterService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(req RegisterRequest) (*RegisterResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	existing, err := s.repo.FindByEmail(req.Email)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("email already registered")
	}

	// 2. Hash password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// 3. Simpan user baru
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// 4. Ambil role default dan assign ke user
	defaultRole, err := s.repo.FindDefaultRole()
	if err == nil && defaultRole.ID != 0 {
		userRole := &models.UserRole{
			UserID: user.ID,
			RoleID: defaultRole.ID,
		}
		_ = s.repo.CreateUserRole(userRole)
	}

	// 5. Return response
	return &RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}
