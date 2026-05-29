package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =========================================
// Helper
// =========================================

func setupPermissionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
	); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// seedUserWithPermissions membuat user lengkap dengan role dan permissions
func seedUserWithPermissions(db *gorm.DB, email string, permissions []string) (*models.User, string) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Buat user
	user := &models.User{Email: email, Password: "hashed"}
	db.Create(user)

	// Buat role
	role := &models.Role{Name: "testrole-" + email, SystemFunction: "test"}
	db.Create(role)

	// Buat permissions
	for _, perm := range permissions {
		db.Create(&models.Permission{RoleID: role.ID, Name: perm})
	}

	// Assign role ke user
	db.Create(&models.UserRole{UserID: user.ID, RoleID: role.ID})

	// Generate token
	token, _ := util.GenerateToken(user.ID)
	return user, token
}

func setupPermissionApp(db *gorm.DB, permissions ...string) *fiber.App {
	app := fiber.New()
	app.Get("/test", middleware.HasPermission(db, permissions...), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user_id": c.Locals("user_id")})
	})
	return app
}

// =========================================
// Test HasPermission — Single Permission
// =========================================

func TestHasPermission_SinglePermission_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)
	_, token := seedUserWithPermissions(db, "user1@example.com", []string{"role:read"})

	app := setupPermissionApp(db, "role:read")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHasPermission_SinglePermission_NotOwned(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)
	// User hanya punya role:read, tapi endpoint butuh role:delete
	_, token := seedUserWithPermissions(db, "user2@example.com", []string{"role:read"})

	app := setupPermissionApp(db, "role:delete")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

// =========================================
// Test HasPermission — Multiple Permissions
// =========================================

func TestHasPermission_MultiplePermissions_AllOwned_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)
	// User punya semua permission yang dibutuhkan
	_, token := seedUserWithPermissions(db, "user3@example.com", []string{
		"role:read", "role:create", "role:update",
	})

	app := setupPermissionApp(db, "role:read", "role:create", "role:update")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHasPermission_MultiplePermissions_SomeNotOwned(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)
	// User hanya punya role:read dan role:create, tapi endpoint butuh role:delete juga
	_, token := seedUserWithPermissions(db, "user4@example.com", []string{
		"role:read", "role:create",
	})

	app := setupPermissionApp(db, "role:read", "role:create", "role:delete")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

func TestHasPermission_MultiplePermissions_NoneOwned(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)
	// User tidak punya permission apapun
	_, token := seedUserWithPermissions(db, "user5@example.com", []string{})

	app := setupPermissionApp(db, "role:read", "role:create")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, res.StatusCode)
}

// =========================================
// Test HasPermission — Auth Error
// =========================================

func TestHasPermission_WithoutToken_ShouldReturn401(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)

	app := setupPermissionApp(db, "role:read")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestHasPermission_InvalidToken_ShouldReturn401(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)

	app := setupPermissionApp(db, "role:read")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token-tidak-valid")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestHasPermission_WrongSecret_ShouldReturn401(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret-A")
	db := setupPermissionTestDB(t)
	_, token := seedUserWithPermissions(db, "user6@example.com", []string{"role:read"})

	// Middleware pakai secret berbeda
	os.Setenv("JWT_SECRET", "secret-B")
	app := setupPermissionApp(db, "role:read")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// =========================================
// Test HasPermission — UserID tersimpan di Locals
// =========================================

func TestHasPermission_SetsUserIDInLocals(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	db := setupPermissionTestDB(t)
	user, token := seedUserWithPermissions(db, "user7@example.com", []string{"role:read"})

	app := fiber.New()
	app.Get("/test", middleware.HasPermission(db, "role:read"), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		return c.JSON(fiber.Map{"user_id": userID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	_ = user
}
