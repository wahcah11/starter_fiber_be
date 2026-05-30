package crud_permission_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/permission/crud_permission/mocks"

	"github.com/stretchr/testify/assert"
)

// =========================================
// FindAll
// =========================================

func TestRepository_FindAll_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := []models.Permission{
		{RoleID: 1, Name: "role:read"},
		{RoleID: 1, Name: "role:create"},
		{RoleID: 2, Name: "user:read"},
	}
	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(expected, int64(3), nil)

	result, total, err := mockRepo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 3, len(result))
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_FilterByRoleID(t *testing.T) {
	mockRepo := new(mocks.Repository)

	roleID := uint(1)
	expected := []models.Permission{
		{RoleID: 1, Name: "role:read"},
		{RoleID: 1, Name: "role:create"},
	}
	mockRepo.On("FindAll", 1, 10, &roleID).Return(expected, int64(2), nil)

	result, total, err := mockRepo.FindAll(1, 10, &roleID)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.Permission{}, int64(0), nil)

	result, total, err := mockRepo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(result))
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(nil, int64(0), errors.New("db error"))

	result, total, err := mockRepo.FindAll(1, 10, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

// =========================================
// FindByID
// =========================================

func TestRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := &models.Permission{RoleID: 1, Name: "role:read"}
	mockRepo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := mockRepo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "role:read", result.Name)
	assert.Equal(t, uint(1), result.RoleID)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindByID(99)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// =========================================
// FindByRoleIDAndName
// =========================================

func TestRepository_FindByRoleIDAndName_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := &models.Permission{RoleID: 1, Name: "role:read"}
	mockRepo.On("FindByRoleIDAndName", uint(1), "role:read").Return(expected, nil)

	result, err := mockRepo.FindByRoleIDAndName(1, "role:read")

	assert.NoError(t, err)
	assert.Equal(t, "role:read", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindByRoleIDAndName_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByRoleIDAndName", uint(1), "role:delete").Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindByRoleIDAndName(1, "role:delete")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// =========================================
// Create
// =========================================

func TestRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	permission := &models.Permission{RoleID: 1, Name: "role:delete"}
	mockRepo.On("Create", permission).Return(nil)

	err := mockRepo.Create(permission)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Create_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	permission := &models.Permission{RoleID: 1, Name: "role:delete"}
	mockRepo.On("Create", permission).Return(errors.New("duplicate entry"))

	err := mockRepo.Create(permission)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// =========================================
// Update
// =========================================

func TestRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	permission := &models.Permission{RoleID: 1, Name: "role:delete-updated"}
	mockRepo.On("Update", permission).Return(nil)

	err := mockRepo.Update(permission)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Update_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	permission := &models.Permission{RoleID: 1, Name: "role:delete-updated"}
	mockRepo.On("Update", permission).Return(errors.New("db error"))

	err := mockRepo.Update(permission)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// =========================================
// Delete
// =========================================

func TestRepository_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := mockRepo.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("Delete", uint(99)).Return(errors.New("record not found"))

	err := mockRepo.Delete(99)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
