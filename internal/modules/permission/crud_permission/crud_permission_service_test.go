package crud_permission_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/permission/crud_permission"
	"starter-wahcah-be/internal/modules/permission/crud_permission/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =========================================
// GetAll
// =========================================

func TestService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.Permission{
		{RoleID: 1, Name: "role:read"},
		{RoleID: 1, Name: "role:create"},
		{RoleID: 2, Name: "user:read"},
	}, int64(3), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Permissions retrieved successfully", res.Message)
	assert.Equal(t, 3, len(res.Data))
	assert.Equal(t, int64(3), res.Pagination.Total)
	assert.Equal(t, 1, res.Pagination.Page)
	assert.Equal(t, 10, res.Pagination.Limit)
	assert.Equal(t, 1, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_FilterByRoleID(t *testing.T) {
	mockRepo := new(mocks.Repository)

	roleID := uint(1)
	mockRepo.On("FindAll", 1, 10, &roleID).Return([]models.Permission{
		{RoleID: 1, Name: "role:read"},
		{RoleID: 1, Name: "role:create"},
	}, int64(2), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 10, &roleID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Data))
	assert.Equal(t, int64(2), res.Pagination.Total)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.Permission{}, int64(0), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Data))
	assert.Equal(t, int64(0), res.Pagination.Total)
	assert.Equal(t, 0, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_InvalidPage_DefaultsToOne(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.Permission{}, int64(0), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(0, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Pagination.Page)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_InvalidLimit_DefaultsToTen(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.Permission{}, int64(0), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 0, nil)

	assert.NoError(t, err)
	assert.Equal(t, 10, res.Pagination.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_LimitExceedsMax_DefaultsToTen(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.Permission{}, int64(0), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 999, nil)

	assert.NoError(t, err)
	assert.Equal(t, 10, res.Pagination.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(nil, int64(0), errors.New("db error"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to retrieve permissions", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_PaginationCalculation(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Total 25 data dengan limit 10 = 3 halaman
	permissions := make([]models.Permission, 10)
	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(permissions, int64(25), nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(25), res.Pagination.Total)
	assert.Equal(t, 3, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

// =========================================
// GetByID
// =========================================

func TestService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.Permission{RoleID: 1, Name: "role:read"}, nil,
	)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Permission retrieved successfully", res.Message)
	assert.Equal(t, "role:read", res.Data.Name)
	assert.Equal(t, uint(1), res.Data.RoleID)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.GetByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "permission not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Create
// =========================================

func TestService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByRoleIDAndName", uint(1), "role:delete").
		Return(nil, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(p *models.Permission) bool {
		return p.RoleID == 1 && p.Name == "role:delete"
	})).Return(nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Create(crud_permission.CreatePermissionRequest{
		RoleID: 1,
		Name:   "role:delete",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Permission created successfully", res.Message)
	assert.Equal(t, "role:delete", res.Data.Name)
	assert.Equal(t, uint(1), res.Data.RoleID)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_AlreadyExists(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Permission{RoleID: 1, Name: "role:read"}
	existing.ID = 1
	mockRepo.On("FindByRoleIDAndName", uint(1), "role:read").Return(existing, nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Create(crud_permission.CreatePermissionRequest{
		RoleID: 1,
		Name:   "role:read",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "permission already exists for this role", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Create_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByRoleIDAndName", uint(1), "role:delete").
		Return(nil, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(p *models.Permission) bool {
		return p.RoleID == 1 && p.Name == "role:delete"
	})).Return(errors.New("db error"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Create(crud_permission.CreatePermissionRequest{
		RoleID: 1,
		Name:   "role:delete",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to create permission", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Update
// =========================================

func TestService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Permission{RoleID: 1, Name: "role:read"}
	existing.ID = 1
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByRoleIDAndName", uint(1), "role:read-updated").
		Return(nil, errors.New("record not found"))
	mockRepo.On("Update", mock.MatchedBy(func(p *models.Permission) bool {
		return p.Name == "role:read-updated" && p.RoleID == 1
	})).Return(nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Update(1, crud_permission.UpdatePermissionRequest{
		RoleID: 1,
		Name:   "role:read-updated",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Permission updated successfully", res.Message)
	assert.Equal(t, "role:read-updated", res.Data.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Update(99, crud_permission.UpdatePermissionRequest{
		RoleID: 1,
		Name:   "role:read",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "permission not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Update_AlreadyExistsForOther(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Permission{RoleID: 1, Name: "role:read"}
	existing.ID = 1
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	// Kombinasi role_id + name sudah dipakai permission lain (ID 2)
	other := &models.Permission{RoleID: 1, Name: "role:create"}
	other.ID = 2
	mockRepo.On("FindByRoleIDAndName", uint(1), "role:create").Return(other, nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Update(1, crud_permission.UpdatePermissionRequest{
		RoleID: 1,
		Name:   "role:create",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "permission already exists for this role", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Update_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Permission{RoleID: 1, Name: "role:read"}
	existing.ID = 1
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByRoleIDAndName", uint(1), "role:read-updated").
		Return(nil, errors.New("record not found"))
	mockRepo.On("Update", mock.MatchedBy(func(p *models.Permission) bool {
		return p.Name == "role:read-updated"
	})).Return(errors.New("db error"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	res, err := svc.Update(1, crud_permission.UpdatePermissionRequest{
		RoleID: 1,
		Name:   "role:read-updated",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to update permission", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Delete
// =========================================

func TestService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.Permission{RoleID: 1, Name: "role:read"}, nil,
	)
	mockRepo.On("Delete", uint(1)).Return(nil)

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	err := svc.Delete(99)

	assert.Error(t, err)
	assert.Equal(t, "permission not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.Permission{RoleID: 1, Name: "role:read"}, nil,
	)
	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	svc := crud_permission.NewCrudPermissionService(mockRepo)
	err := svc.Delete(1)

	assert.Error(t, err)
	assert.Equal(t, "failed to delete permission", err.Error())
	mockRepo.AssertExpectations(t)
}
