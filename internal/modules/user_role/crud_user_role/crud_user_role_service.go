package crud_user_role

import (
	"errors"
	"math"

	"starter-wahcah-be/internal/models"
)

type Service interface {
	GetAll(page, limit int, userID *uint) (*UserRoleListResponse, error)
	GetByID(id uint) (*UserRoleSingleResponse, error)
	Create(req CreateUserRoleRequest) (*UserRoleSingleResponse, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewCrudUserRoleService(repo Repository) Service {
	return &service{repo}
}

// toResponse mengkonversi models.UserRole ke UserRoleResponse
func toResponse(ur models.UserRole) UserRoleResponse {
	return UserRoleResponse{
		ID:     ur.ID,
		UserID: ur.UserID,
		RoleID: ur.RoleID,
	}
}

func (s *service) GetAll(page, limit int, userID *uint) (*UserRoleListResponse, error) {
	// Validasi page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	userRoles, total, err := s.repo.FindAll(page, limit, userID)
	if err != nil {
		return nil, errors.New("failed to retrieve user roles")
	}

	// Konversi ke response
	data := make([]UserRoleResponse, len(userRoles))
	for i, ur := range userRoles {
		data[i] = toResponse(ur)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &UserRoleListResponse{
		Message: "User roles retrieved successfully",
		Data:    data,
		Pagination: PaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (s *service) GetByID(id uint) (*UserRoleSingleResponse, error) {
	userRole, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user role not found")
	}

	return &UserRoleSingleResponse{
		Message: "User role retrieved successfully",
		Data:    toResponse(*userRole),
	}, nil
}

func (s *service) Create(req CreateUserRoleRequest) (*UserRoleSingleResponse, error) {
	// Cek duplikat — kombinasi user_id + role_id harus unik
	existing, err := s.repo.FindByUserIDAndRoleID(req.UserID, req.RoleID)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("role already assigned to this user")
	}

	userRole := &models.UserRole{
		UserID: req.UserID,
		RoleID: req.RoleID,
	}

	if err := s.repo.Create(userRole); err != nil {
		return nil, errors.New("failed to assign role to user")
	}

	return &UserRoleSingleResponse{
		Message: "Role assigned to user successfully",
		Data:    toResponse(*userRole),
	}, nil
}

func (s *service) Delete(id uint) error {
	// Cek apakah user_role ada
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user role not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return errors.New("failed to remove role from user")
	}

	return nil
}
