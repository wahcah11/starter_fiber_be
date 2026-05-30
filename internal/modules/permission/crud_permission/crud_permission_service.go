package crud_permission

import (
	"errors"
	"math"

	"starter-wahcah-be/internal/models"
)

type Service interface {
	GetAll(page, limit int, roleID *uint) (*PermissionListResponse, error)
	GetByID(id uint) (*PermissionSingleResponse, error)
	Create(req CreatePermissionRequest) (*PermissionSingleResponse, error)
	Update(id uint, req UpdatePermissionRequest) (*PermissionSingleResponse, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewCrudPermissionService(repo Repository) Service {
	return &service{repo}
}

// toResponse mengkonversi models.Permission ke PermissionResponse
func toResponse(p models.Permission) PermissionResponse {
	return PermissionResponse{
		ID:     p.ID,
		RoleID: p.RoleID,
		Name:   p.Name,
	}
}

func (s *service) GetAll(page, limit int, roleID *uint) (*PermissionListResponse, error) {
	// Validasi page dan limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	permissions, total, err := s.repo.FindAll(page, limit, roleID)
	if err != nil {
		return nil, errors.New("failed to retrieve permissions")
	}

	// Konversi ke response
	data := make([]PermissionResponse, len(permissions))
	for i, p := range permissions {
		data[i] = toResponse(p)
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &PermissionListResponse{
		Message: "Permissions retrieved successfully",
		Data:    data,
		Pagination: PaginationMeta{
			Page:      page,
			Limit:     limit,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

func (s *service) GetByID(id uint) (*PermissionSingleResponse, error) {
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}

	return &PermissionSingleResponse{
		Message: "Permission retrieved successfully",
		Data:    toResponse(*permission),
	}, nil
}

func (s *service) Create(req CreatePermissionRequest) (*PermissionSingleResponse, error) {
	// Cek duplikat — kombinasi role_id + name harus unik
	existing, err := s.repo.FindByRoleIDAndName(req.RoleID, req.Name)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("permission already exists for this role")
	}

	permission := &models.Permission{
		RoleID: req.RoleID,
		Name:   req.Name,
	}

	if err := s.repo.Create(permission); err != nil {
		return nil, errors.New("failed to create permission")
	}

	return &PermissionSingleResponse{
		Message: "Permission created successfully",
		Data:    toResponse(*permission),
	}, nil
}

func (s *service) Update(id uint, req UpdatePermissionRequest) (*PermissionSingleResponse, error) {
	// Cek apakah permission ada
	permission, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("permission not found")
	}

	// Cek duplikat — kombinasi role_id + name baru tidak boleh dipakai permission lain
	existing, err := s.repo.FindByRoleIDAndName(req.RoleID, req.Name)
	if err == nil && existing.ID != 0 && existing.ID != id {
		return nil, errors.New("permission already exists for this role")
	}

	permission.RoleID = req.RoleID
	permission.Name = req.Name

	if err := s.repo.Update(permission); err != nil {
		return nil, errors.New("failed to update permission")
	}

	return &PermissionSingleResponse{
		Message: "Permission updated successfully",
		Data:    toResponse(*permission),
	}, nil
}

func (s *service) Delete(id uint) error {
	// Cek apakah permission ada
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("permission not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return errors.New("failed to delete permission")
	}

	return nil
}
