package crud_user_role_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/user_role/crud_user_role/mocks"

	"github.com/stretchr/testify/assert"
)

// =========================================
// FindAll
// =========================================

func TestRepository_FindAll_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := []models.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 1, RoleID: 2},
		{UserID: 2, RoleID: 1},
	}
	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(expected, int64(3), nil)

	result, total, err := mockRepo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 3, len(result))
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_FilterByUserID(t *testing.T) {
	mockRepo := new(mocks.Repository)

	userID := uint(1)
	expected := []models.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 1, RoleID: 2},
	}
	mockRepo.On("FindAll", 1, 10, &userID).Return(expected, int64(2), nil)

	result, total, err := mockRepo.FindAll(1, 10, &userID)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))
	for _, ur := range result {
		assert.Equal(t, uint(1), ur.UserID)
	}
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.UserRole{}, int64(0), nil)

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

	expected := &models.UserRole{UserID: 1, RoleID: 1}
	mockRepo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := mockRepo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
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
// FindByUserIDAndRoleID
// =========================================

func TestRepository_FindByUserIDAndRoleID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := &models.UserRole{UserID: 1, RoleID: 2}
	mockRepo.On("FindByUserIDAndRoleID", uint(1), uint(2)).Return(expected, nil)

	result, err := mockRepo.FindByUserIDAndRoleID(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, uint(2), result.RoleID)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindByUserIDAndRoleID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByUserIDAndRoleID", uint(1), uint(99)).Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindByUserIDAndRoleID(1, 99)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// =========================================
// Create
// =========================================

func TestRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	userRole := &models.UserRole{UserID: 1, RoleID: 3}
	mockRepo.On("Create", userRole).Return(nil)

	err := mockRepo.Create(userRole)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Create_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	userRole := &models.UserRole{UserID: 1, RoleID: 3}
	mockRepo.On("Create", userRole).Return(errors.New("duplicate entry"))

	err := mockRepo.Create(userRole)

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

func TestRepository_Delete_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("Delete", uint(99)).Return(errors.New("record not found"))

	err := mockRepo.Delete(99)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
