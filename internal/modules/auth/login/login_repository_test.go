package login_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/login/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FindByEmail_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	expectedUser := &models.User{
		Email:    "member@example.com",
		Password: "hashedpassword",
	}

	// Simulasi: FindByEmail dipanggil dengan email tersebut dan return user
	mockRepo.On("FindByEmail", "member@example.com").Return(expectedUser, nil)

	result, err := mockRepo.FindByEmail("member@example.com")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestRepository_FindByEmail_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Simulasi: email tidak ditemukan
	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("record not found"))

	result, err := mockRepo.FindByEmail("notfound@example.com")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}
