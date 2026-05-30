package crud_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/role/crud/mocks"

	"github.com/stretchr/testify/assert"
)

// =========================================
// FindAll
// =========================================

func TestRepository_FindAll_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := []models.Role{
		{Name: "superadmin", SystemFunction: "full_access", IsDefault: false},
		{Name: "member", SystemFunction: "basic_access", IsDefault: true},
	}
	mockRepo.On("FindAll", 1, 10).Return(expected, int64(2), nil)

	result, total, err := mockRepo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "superadmin", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return([]models.Role{}, int64(0), nil)

	result, total, err := mockRepo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(result))
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindAll_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10).Return(nil, int64(0), errors.New("db error"))

	result, total, err := mockRepo.FindAll(1, 10)

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

	expected := &models.Role{Name: "superadmin", SystemFunction: "full_access"}
	mockRepo.On("FindByID", uint(1)).Return(expected, nil)

	result, err := mockRepo.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, "superadmin", result.Name)
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
// FindByName
// =========================================

func TestRepository_FindByName_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := &models.Role{Name: "member", SystemFunction: "basic_access"}
	mockRepo.On("FindByName", "member").Return(expected, nil)

	result, err := mockRepo.FindByName("member")

	assert.NoError(t, err)
	assert.Equal(t, "member", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindByName_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByName", "unknown").Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindByName("unknown")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// =========================================
// Create
// =========================================

func TestRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	role := &models.Role{Name: "editor", SystemFunction: "edit_access"}
	mockRepo.On("Create", role).Return(nil)

	err := mockRepo.Create(role)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Create_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	role := &models.Role{Name: "editor", SystemFunction: "edit_access"}
	mockRepo.On("Create", role).Return(errors.New("duplicate entry"))

	err := mockRepo.Create(role)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// =========================================
// Update
// =========================================

func TestRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	role := &models.Role{Name: "editor-updated", SystemFunction: "edit_access"}
	mockRepo.On("Update", role).Return(nil)

	err := mockRepo.Update(role)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Update_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	role := &models.Role{Name: "editor-updated", SystemFunction: "edit_access"}
	mockRepo.On("Update", role).Return(errors.New("db error"))

	err := mockRepo.Update(role)

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
