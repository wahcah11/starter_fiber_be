package register_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/register"
	"starter-wahcah-be/internal/modules/auth/register/mocks"

	"github.com/stretchr/testify/assert"
)

func TestService_Register_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Email belum terdaftar
	mockRepo.On("FindByEmail", "newuser@example.com").
		Return(&models.User{}, errors.New("record not found"))

	// Create user berhasil
	mockRepo.On("Create", mock_AnyUser()).
		Return(nil)

	// Ada role default
	mockRepo.On("FindDefaultRole").
		Return(&models.Role{Name: "member", IsDefault: true}, nil)

	// Assign role berhasil
	mockRepo.On("CreateUserRole", mock_AnyUserRole()).
		Return(nil)

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

	// Email sudah terdaftar — FindByEmail return user yang ada
	mockRepo.On("FindByEmail", "existing@example.com").
		Return(&models.User{Email: "existing@example.com"}, nil)

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

	mockRepo.On("Create", mock_AnyUser()).
		Return(errors.New("db error"))

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
	mockRepo.On("Create", mock_AnyUser()).
		Return(nil)

	// Tidak ada role default — tetap lanjut, tidak return error
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

// =========================================
// Helpers — mock matcher untuk pointer struct
// =========================================

// mock_AnyUser mengembalikan matcher untuk *models.User apapun
// karena nilai pointer tidak bisa diprediksi sebelum Create dipanggil
func mock_AnyUser() interface{} {
	return mock_matcherFunc(func(arg interface{}) bool {
		_, ok := arg.(*models.User)
		return ok
	})
}

// mock_AnyUserRole mengembalikan matcher untuk *models.UserRole apapun
func mock_AnyUserRole() interface{} {
	return mock_matcherFunc(func(arg interface{}) bool {
		_, ok := arg.(*models.UserRole)
		return ok
	})
}

// mock_matcherFunc adalah custom matcher untuk testify mock
type mock_matcherFunc func(interface{}) bool

func (f mock_matcherFunc) Matches(arg interface{}) bool {
	return f(arg)
}

func (f mock_matcherFunc) String() string {
	return "custom matcher"
}
