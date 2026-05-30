package crud

import (
	"errors"
	"math"

	"starter-wahcah-be/internal/models"
)

type Service interface {
	GetAll(page, limit int) (*RoleListResponse, error)
	GetByID(id uint) (*RoleSingleResponse, error)
	Create(req CreateRoleRequest) (*RoleSingleResponse, error)
	Update(id uint, req UpdateRoleRequest) (*RoleSingleResponse, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewCrudService(repo Repository) Service {
	return &service{repo}
}

// toResponse mengkonversi models.Role ke RoleResponse
func toResponse(role models.Role) RoleResponse {
	return RoleResponse{
		ID:             role.ID,
		Name:           role.Name,
		SystemFunction: role.SystemFunction,
		IsDefault:      role.IsDefault,
	}
}

func (s *service) GetAll(page, limit int) (*RoleListResponse, error) {
	// Validasi page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	roles, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, errors.New("failed to retrieve roles")
	}

	// Konversi ke response
	data := make([]RoleResponse, len(roles))
	for i, role := range roles {
		data[i] = toResponse(role)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &RoleListResponse{
		Message: "Roles retrieved successfully",
		Data:    data,
		Pagination: PaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (s *service) GetByID(id uint) (*RoleSingleResponse, error) {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	return &RoleSingleResponse{
		Message: "Role retrieved successfully",
		Data:    toResponse(*role),
	}, nil
}

func (s *service) Create(req CreateRoleRequest) (*RoleSingleResponse, error) {
	// Cek apakah nama role sudah ada
	existing, err := s.repo.FindByName(req.Name)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("role name already exists")
	}

	role := &models.Role{
		Name:           req.Name,
		SystemFunction: req.SystemFunction,
		IsDefault:      req.IsDefault,
	}

	if err := s.repo.Create(role); err != nil {
		return nil, errors.New("failed to create role")
	}

	return &RoleSingleResponse{
		Message: "Role created successfully",
		Data:    toResponse(*role),
	}, nil
}

func (s *service) Update(id uint, req UpdateRoleRequest) (*RoleSingleResponse, error) {
	// Cek apakah role ada
	role, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Cek apakah nama baru sudah dipakai role lain
	existing, err := s.repo.FindByName(req.Name)
	if err == nil && existing.ID != 0 && existing.ID != id {
		return nil, errors.New("role name already exists")
	}

	role.Name = req.Name
	role.SystemFunction = req.SystemFunction
	role.IsDefault = req.IsDefault

	if err := s.repo.Update(role); err != nil {
		return nil, errors.New("failed to update role")
	}

	return &RoleSingleResponse{
		Message: "Role updated successfully",
		Data:    toResponse(*role),
	}, nil
}

func (s *service) Delete(id uint) error {
	// Cek apakah role ada
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return errors.New("failed to delete role")
	}

	return nil
}
