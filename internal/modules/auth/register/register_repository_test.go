package register_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/register/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FindByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := &models.User{Email: "user@example.com", Password: "hashed"}
	mockRepo.On("FindByEmail", "user@example.com").Return(expected, nil)

	result, err := mockRepo.FindByEmail("user@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expected.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindByEmail_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByEmail", "ghost@example.com").Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindByEmail("ghost@example.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	user := &models.User{Email: "newuser@example.com", Password: "hashed"}
	mockRepo.On("Create", user).Return(nil)

	err := mockRepo.Create(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_Create_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	user := &models.User{Email: "newuser@example.com", Password: "hashed"}
	mockRepo.On("Create", user).Return(errors.New("duplicate entry"))

	err := mockRepo.Create(user)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindDefaultRole_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expected := &models.Role{Name: "member", IsDefault: true}
	mockRepo.On("FindDefaultRole").Return(expected, nil)

	result, err := mockRepo.FindDefaultRole()

	assert.NoError(t, err)
	assert.Equal(t, "member", result.Name)
	assert.True(t, result.IsDefault)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindDefaultRole_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindDefaultRole").Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindDefaultRole()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRepository_CreateUserRole_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	userRole := &models.UserRole{UserID: 1, RoleID: 1}
	mockRepo.On("CreateUserRole", userRole).Return(nil)

	err := mockRepo.CreateUserRole(userRole)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRepository_CreateUserRole_Failed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	userRole := &models.UserRole{UserID: 1, RoleID: 1}
	mockRepo.On("CreateUserRole", userRole).Return(errors.New("db error"))

	err := mockRepo.CreateUserRole(userRole)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
