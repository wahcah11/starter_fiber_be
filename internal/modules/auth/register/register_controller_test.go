package register_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"starter-wahcah-be/internal/modules/auth/register"
	"starter-wahcah-be/internal/modules/auth/register/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// =========================================
// Helper
// =========================================

func setupRegisterControllerApp(svc register.Service) *fiber.App {
	app := fiber.New()
	ctrl := register.NewRegisterController(svc)
	app.Post("/auth/register", ctrl.Register)
	return app
}

// =========================================
// Test Controller
// =========================================

func TestController_Register_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Register", register.RegisterRequest{
		Email:    "newuser@example.com",
		Password: "secret123",
	}).Return(&register.RegisterResponse{
		ID:    1,
		Email: "newuser@example.com",
	}, nil)

	app := setupRegisterControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "newuser@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "newuser@example.com", data["email"])
	mockSvc.AssertExpectations(t)
}

func TestController_Register_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupRegisterControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Invalid JSON", response["error"])
}

func TestController_Register_ValidationError_InvalidEmail(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupRegisterControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "bukan-email",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["validation"])
}

func TestController_Register_ValidationError_PasswordTooShort(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupRegisterControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "newuser@example.com",
		"password": "123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["validation"])
}

func TestController_Register_ValidationError_EmptyFields(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupRegisterControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Register_EmailAlreadyRegistered(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Register", register.RegisterRequest{
		Email:    "existing@example.com",
		Password: "secret123",
	}).Return(nil, assert.AnError)

	app := setupRegisterControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "existing@example.com",
		"password": "secret123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["error"])
	mockSvc.AssertExpectations(t)
}
