package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// =========================================
// Helper
// =========================================

// setupProtectedApp membuat app fiber dengan middleware Protected
func setupProtectedApp() *fiber.App {
	app := fiber.New()
	app.Get("/protected", middleware.Protected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"user_id": userID})
	})
	return app
}

// setupOptionalApp membuat app fiber dengan middleware OptionalAuth
func setupOptionalApp() *fiber.App {
	app := fiber.New()
	app.Get("/optional", middleware.OptionalAuth(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"user_id": userID})
	})
	return app
}

// generateTestToken membuat JWT token valid untuk test
func generateTestToken(userID uint) string {
	os.Setenv("JWT_SECRET", "test-secret-key")
	token, _ := util.GenerateToken(userID)
	return token
}

// =========================================
// Protected Middleware
// =========================================

func TestProtected_WithValidToken_ShouldPass(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupProtectedApp()

	token := generateTestToken(1)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestProtected_WithValidToken_ShouldSetUserID(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	app := fiber.New()
	app.Get("/protected", middleware.Protected(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"user_id": userID})
	})

	token := generateTestToken(42)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestProtected_WithoutToken_ShouldReturn401(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupProtectedApp()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	// Tidak set Authorization header

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestProtected_WithInvalidToken_ShouldReturn401(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupProtectedApp()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer token-tidak-valid")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestProtected_WithWrongSecret_ShouldReturn401(t *testing.T) {
	// Generate token dengan secret A
	os.Setenv("JWT_SECRET", "secret-A")
	token := generateTestToken(1)

	// Tapi middleware menggunakan secret B
	os.Setenv("JWT_SECRET", "secret-B")
	app := setupProtectedApp()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestProtected_WithEmptyBearerToken_ShouldReturn401(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupProtectedApp()

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// =========================================
// OptionalAuth Middleware
// =========================================

func TestOptionalAuth_WithValidToken_ShouldSetUserID(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupOptionalApp()

	token := generateTestToken(7)
	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	// Harus tetap 200, bukan 401
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestOptionalAuth_WithoutToken_ShouldStillPass(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupOptionalApp()

	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	// Tidak set Authorization header sama sekali

	res, err := app.Test(req)

	assert.NoError(t, err)
	// Harus 200, bukan 401 — ini bedanya dengan Protected
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestOptionalAuth_WithInvalidToken_ShouldStillPass(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	app := setupOptionalApp()

	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	req.Header.Set("Authorization", "Bearer token-tidak-valid")

	res, err := app.Test(req)

	assert.NoError(t, err)
	// Token invalid tapi tetap lanjut, bukan 401
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestOptionalAuth_WithWrongSecret_ShouldStillPass(t *testing.T) {
	// Generate token dengan secret A
	os.Setenv("JWT_SECRET", "secret-A")
	token := generateTestToken(1)

	// Middleware pakai secret B — token invalid tapi tetap lanjut
	os.Setenv("JWT_SECRET", "secret-B")
	app := setupOptionalApp()

	req := httptest.NewRequest(http.MethodGet, "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
