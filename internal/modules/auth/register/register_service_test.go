package register_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/register"
	"starter-wahcah-be/internal/modules/auth/register/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Email belum terdaftar
	mockRepo.On("FindByEmail", "newuser@example.com").
		Return(&models.User{}, errors.New("record not found"))

	// Create user berhasil — set ID via Run agar CreateUserRole bisa dipanggil
	mockRepo.On("Create", mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "newuser@example.com"
	})).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		user.ID = 1 // simulasi ID yang di-set DB setelah insert
	}).Return(nil)

	// Ada role default dengan ID
	defaultRole := &models.Role{Name: "member", IsDefault: true}
	defaultRole.ID = 1
	mockRepo.On("FindDefaultRole").Return(defaultRole, nil)

	// Assign role dipanggil karena user.ID dan defaultRole.ID keduanya != 0
	mockRepo.On("CreateUserRole", mock.MatchedBy(func(ur *models.UserRole) bool {
		return ur.UserID == 1 && ur.RoleID == 1
	})).Return(nil)

	svc := register.NewRegisterService(mockRepo)
	res, err := svc.Register(register.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "secret123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "newuser@example.com", res.Email)
	mockRepo.AssertExpectations(t)
}

func TestService_Register_EmailAlreadyRegistered(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.User{Email: "existing@example.com"}
	existing.ID = 1
	mockRepo.On("FindByEmail", "existing@example.com").Return(existing, nil)

	svc := register.NewRegisterService(mockRepo)
	res, err := svc.Register(register.RegisterRequest{
		Email:    "existing@example.com",
		Password: "secret123",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "email already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Register_CreateFailed(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByEmail", "newuser@example.com").
		Return(&models.User{}, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "newuser@example.com"
	})).Return(errors.New("db error"))

	svc := register.NewRegisterService(mockRepo)
	res, err := svc.Register(register.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "secret123",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to create user", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Register_NoDefaultRole_StillSuccess(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Email belum terdaftar
	mockRepo.On("FindByEmail", "newuser@example.com").
		Return(&models.User{}, errors.New("record not found"))

	// Create user berhasil
	mockRepo.On("Create", mock.MatchedBy(func(u *models.User) bool {
		return u.Email == "newuser@example.com"
	})).Return(nil)

	// Tidak ada role default — tetap lanjut
	mockRepo.On("FindDefaultRole").
		Return(&models.Role{}, errors.New("record not found"))

	svc := register.NewRegisterService(mockRepo)
	res, err := svc.Register(register.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "secret123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "newuser@example.com", res.Email)
	mockRepo.AssertExpectations(t)
}
