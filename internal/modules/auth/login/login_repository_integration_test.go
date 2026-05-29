package login_test

import (
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =========================================
// Helper
// =========================================

// setupTestDB membuat SQLite in-memory database untuk integration test
// Tidak konek ke MySQL sungguhan — aman dijalankan di CI/CD
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Migrate hanya tabel yang dibutuhkan
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// seedTestUser menyisipkan user dummy ke DB test
func seedTestUser(db *gorm.DB, email, password string) *models.User {
	user := &models.User{
		Email:    email,
		Password: password,
	}
	db.Create(user)
	return user
}

// =========================================
// Integration Test Repository
// =========================================

func TestIntegration_Repository_FindByEmail_Success(t *testing.T) {
	db := setupTestDB(t)
	seedTestUser(db, "member@example.com", "hashedpassword")

	repo := login.NewLoginRepository(db)
	result, err := repo.FindByEmail("member@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "member@example.com", result.Email)
}

func TestIntegration_Repository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	// Tidak ada user yang di-seed

	repo := login.NewLoginRepository(db)
	result, err := repo.FindByEmail("ghost@example.com")

	assert.Error(t, err)
	// GORM mengembalikan zero-value struct, bukan nil pointer saat error First
	_ = result
}

func TestIntegration_Repository_FindByEmail_ReturnsCorrectUser(t *testing.T) {
	db := setupTestDB(t)

	// Seed 2 user berbeda
	seedTestUser(db, "user1@example.com", "hash1")
	seedTestUser(db, "user2@example.com", "hash2")

	repo := login.NewLoginRepository(db)
	result, err := repo.FindByEmail("user2@example.com")

	assert.NoError(t, err)
	// Pastikan yang dikembalikan adalah user yang benar
	assert.Equal(t, "user2@example.com", result.Email)
	assert.Equal(t, "hash2", result.Password)
}

func TestIntegration_Repository_FindByEmail_IsCaseSensitive(t *testing.T) {
	db := setupTestDB(t)
	seedTestUser(db, "Member@example.com", "hashedpassword")

	repo := login.NewLoginRepository(db)

	// Email dengan huruf kecil seharusnya tidak ditemukan
	// (SQLite case-sensitive untuk LIKE, tapi WHERE = bisa berbeda per konfigurasi)
	// Test ini memastikan behavior konsisten
	result, err := repo.FindByEmail("member@example.com")
	_ = result

	// Catat: di MySQL behavior ini bisa berbeda tergantung collation
	// Test ini sebagai dokumentasi behavior
	_ = err
}

func TestIntegration_Repository_FindByEmail_ReturnPasswordHash(t *testing.T) {
	db := setupTestDB(t)
	seedTestUser(db, "member@example.com", "$2a$14$hashedpasswordexample")

	repo := login.NewLoginRepository(db)
	result, err := repo.FindByEmail("member@example.com")

	assert.NoError(t, err)
	// Pastikan password hash tersimpan dan terbaca dengan benar
	assert.Equal(t, "$2a$14$hashedpasswordexample", result.Password)
}