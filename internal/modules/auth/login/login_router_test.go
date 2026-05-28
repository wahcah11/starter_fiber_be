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

func setupRouterApp(svc login.Service) *fiber.App {
	app := fiber.New()

	// Simulasi router tanpa DB, inject service langsung
	ctrl := login.NewLoginController(svc)
	auth := app.Group("/auth")
	auth.Post("/login", ctrl.Login)

	return app
}

func TestRouter_POST_AuthLogin_RouteExists(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Authenticate", login.LoginRequest{
		Email:    "member@example.com",
		Password: "member123",
	}).Return(&login.LoginResponse{Token: "mocked-token"}, nil)

	app := setupRouterApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"email":    "member@example.com",
		"password": "member123",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	// Pastikan route terdaftar dan tidak return 404
	assert.NotEqual(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestRouter_GET_AuthLogin_MethodNotAllowed(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupRouterApp(mockSvc)

	// Route hanya menerima POST, GET harus 405
	req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}
