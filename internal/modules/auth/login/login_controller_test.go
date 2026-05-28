package login_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/auth/login/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupControllerApp(svc login.Service) *fiber.App {
	app := fiber.New()
	ctrl := login.NewLoginController(svc)
	app.Post("/auth/login", ctrl.Login)
	return app
}

func TestController_Login_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Authenticate", login.LoginRequest{
		Email:    "member@example.com",
		Password: "member123",
	}).Return(&login.LoginResponse{Token: "mocked-jwt-token"}, nil)

	app := setupControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "member@example.com",
		"password": "member123",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "mocked-jwt-token", data["token"])

	mockSvc.AssertExpectations(t)
}

func TestController_Login_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupControllerApp(mockSvc)

	// Kirim body yang bukan JSON valid
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Login_ValidationError(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupControllerApp(mockSvc)

	// Email tidak valid, password terlalu pendek
	body, _ := json.Marshal(fiber.Map{
		"email":    "bukan-email",
		"password": "123",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Login_Unauthorized(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Authenticate", login.LoginRequest{
		Email:    "member@example.com",
		Password: "wrongpassword",
	}).Return(nil, assert.AnError)

	app := setupControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "member@example.com",
		"password": "wrongpassword",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	mockSvc.AssertExpectations(t)
}
