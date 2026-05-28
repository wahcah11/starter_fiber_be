package login_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/auth/login/mocks"
	"starter-wahcah-be/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestService_Authenticate_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	hashedPassword, _ := util.HashPassword("member123")
	expectedUser := &models.User{
		Email:    "member@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("FindByEmail", "member@example.com").Return(expectedUser, nil)

	svc := login.NewLoginService(mockRepo)
	res, err := svc.Authenticate(login.LoginRequest{
		Email:    "member@example.com",
		Password: "member123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	mockRepo.AssertExpectations(t)
}

func TestService_Authenticate_EmailNotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByEmail", "ghost@example.com").Return(nil, errors.New("record not found"))

	svc := login.NewLoginService(mockRepo)
	res, err := svc.Authenticate(login.LoginRequest{
		Email:    "ghost@example.com",
		Password: "somepassword",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Authenticate_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.Repository)

	hashedPassword, _ := util.HashPassword("correctpassword")
	expectedUser := &models.User{
		Email:    "member@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("FindByEmail", "member@example.com").Return(expectedUser, nil)

	svc := login.NewLoginService(mockRepo)
	res, err := svc.Authenticate(login.LoginRequest{
		Email:    "member@example.com",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}
