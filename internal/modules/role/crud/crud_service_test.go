package crud_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/role/crud"
	"starter-wahcah-be/internal/modules/role/crud/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =========================================
// GetAll
// =========================================

func TestService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return([]models.Role{
		{Name: "superadmin", SystemFunction: "full_access", IsDefault: false},
		{Name: "member", SystemFunction: "basic_access", IsDefault: true},
	}, int64(2), nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Roles retrieved successfully", res.Message)
	assert.Equal(t, 2, len(res.Data))
	assert.Equal(t, int64(2), res.Pagination.Total)
	assert.Equal(t, 1, res.Pagination.Page)
	assert.Equal(t, 10, res.Pagination.Limit)
	assert.Equal(t, 1, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return([]models.Role{}, int64(0), nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Data))
	assert.Equal(t, int64(0), res.Pagination.Total)
	assert.Equal(t, 0, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_InvalidPage_DefaultsToOne(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return([]models.Role{}, int64(0), nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(0, 10)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Pagination.Page)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_InvalidLimit_DefaultsToTen(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return([]models.Role{}, int64(0), nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(1, 0)

	assert.NoError(t, err)
	assert.Equal(t, 10, res.Pagination.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_LimitExceedsMax_DefaultsToTen(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return([]models.Role{}, int64(0), nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(1, 999)

	assert.NoError(t, err)
	assert.Equal(t, 10, res.Pagination.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return(nil, int64(0), errors.New("db error"))

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(1, 10)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to retrieve roles", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_PaginationCalculation(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Total 25 data dengan limit 10 = 3 halaman
	roles := make([]models.Role, 10)
	mockRepo.On("FindAll", 1, 10).Return(roles, int64(25), nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetAll(1, 10)

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
		&models.Role{Name: "superadmin", SystemFunction: "full_access", IsDefault: false}, nil,
	)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Role retrieved successfully", res.Message)
	assert.Equal(t, "superadmin", res.Data.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.GetByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "role not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Create
// =========================================

func TestService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByName", "editor").
		Return(nil, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == "editor" && r.SystemFunction == "edit_access"
	})).Return(nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Create(crud.CreateRoleRequest{
		Name:           "editor",
		SystemFunction: "edit_access",
		IsDefault:      false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Role created successfully", res.Message)
	assert.Equal(t, "editor", res.Data.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_NameAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Role{Name: "member", SystemFunction: "basic_access"}
	existing.ID = 1
	mockRepo.On("FindByName", "member").Return(existing, nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Create(crud.CreateRoleRequest{
		Name:           "member",
		SystemFunction: "basic_access",
		IsDefault:      true,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "role name already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Create_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByName", "editor").
		Return(nil, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == "editor"
	})).Return(errors.New("db error"))

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Create(crud.CreateRoleRequest{
		Name:           "editor",
		SystemFunction: "edit_access",
		IsDefault:      false,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to create role", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Update
// =========================================

func TestService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Role{Name: "editor", SystemFunction: "edit_access"}
	existing.ID = 1
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByName", "editor-updated").
		Return(nil, errors.New("record not found"))
	mockRepo.On("Update", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == "editor-updated"
	})).Return(nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Update(1, crud.UpdateRoleRequest{
		Name:           "editor-updated",
		SystemFunction: "edit_access",
		IsDefault:      false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Role updated successfully", res.Message)
	assert.Equal(t, "editor-updated", res.Data.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_Update_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Update(99, crud.UpdateRoleRequest{
		Name:           "editor",
		SystemFunction: "edit_access",
		IsDefault:      false,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "role not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Update_NameAlreadyUsedByOther(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Role{Name: "editor", SystemFunction: "edit_access"}
	existing.ID = 1
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)

	// Nama "member" sudah dipakai role lain (ID 2)
	otherRole := &models.Role{Name: "member"}
	otherRole.ID = 2
	mockRepo.On("FindByName", "member").Return(otherRole, nil)

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Update(1, crud.UpdateRoleRequest{
		Name:           "member",
		SystemFunction: "edit_access",
		IsDefault:      false,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "role name already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Update_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.Role{Name: "editor", SystemFunction: "edit_access"}
	existing.ID = 1
	mockRepo.On("FindByID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByName", "editor-updated").
		Return(nil, errors.New("record not found"))
	mockRepo.On("Update", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == "editor-updated"
	})).Return(errors.New("db error"))

	svc := crud.NewCrudService(mockRepo)
	res, err := svc.Update(1, crud.UpdateRoleRequest{
		Name:           "editor-updated",
		SystemFunction: "edit_access",
		IsDefault:      false,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to update role", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Delete
// =========================================

func TestService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.Role{Name: "editor"}, nil,
	)
	mockRepo.On("Delete", uint(1)).Return(nil)

	svc := crud.NewCrudService(mockRepo)
	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud.NewCrudService(mockRepo)
	err := svc.Delete(99)

	assert.Error(t, err)
	assert.Equal(t, "role not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.Role{Name: "editor"}, nil,
	)
	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	svc := crud.NewCrudService(mockRepo)
	err := svc.Delete(1)

	assert.Error(t, err)
	assert.Equal(t, "failed to delete role", err.Error())
	mockRepo.AssertExpectations(t)
}
