package register_test

import (
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/auth/register"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =========================================
// Helper
// =========================================

func setupRegisterTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.UserRole{},
	); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func seedRegisterRole(db *gorm.DB, name string, isDefault bool) *models.Role {
	role := &models.Role{Name: name, SystemFunction: "basic_access", IsDefault: isDefault}
	db.Create(role)
	return role
}

// =========================================
// Integration Test Repository
// =========================================

func TestIntegration_Register_FindByEmail_Success(t *testing.T) {
	db := setupRegisterTestDB(t)
	db.Create(&models.User{Email: "existing@example.com", Password: "hashed"})

	repo := register.NewRegisterRepository(db)
	result, err := repo.FindByEmail("existing@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "existing@example.com", result.Email)
}

func TestIntegration_Register_FindByEmail_NotFound(t *testing.T) {
	db := setupRegisterTestDB(t)

	repo := register.NewRegisterRepository(db)
	_, err := repo.FindByEmail("ghost@example.com")

	assert.Error(t, err)
}

func TestIntegration_Register_Create_Success(t *testing.T) {
	db := setupRegisterTestDB(t)

	repo := register.NewRegisterRepository(db)
	user := &models.User{
		Email:    "newuser@example.com",
		Password: "hashedpassword",
	}

	err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestIntegration_Register_Create_DuplicateEmail(t *testing.T) {
	db := setupRegisterTestDB(t)
	db.Create(&models.User{Email: "existing@example.com", Password: "hashed"})

	repo := register.NewRegisterRepository(db)
	duplicate := &models.User{
		Email:    "existing@example.com",
		Password: "hashed2",
	}

	err := repo.Create(duplicate)

	// Harus error karena email unique
	assert.Error(t, err)
}

func TestIntegration_Register_FindDefaultRole_Success(t *testing.T) {
	db := setupRegisterTestDB(t)
	seedRegisterRole(db, "member", true)

	repo := register.NewRegisterRepository(db)
	result, err := repo.FindDefaultRole()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsDefault)
	assert.Equal(t, "member", result.Name)
}

func TestIntegration_Register_FindDefaultRole_NotFound(t *testing.T) {
	db := setupRegisterTestDB(t)
	// Tidak ada role yang di-seed

	repo := register.NewRegisterRepository(db)
	_, err := repo.FindDefaultRole()

	assert.Error(t, err)
}

func TestIntegration_Register_CreateUserRole_Success(t *testing.T) {
	db := setupRegisterTestDB(t)

	user := &models.User{Email: "newuser@example.com", Password: "hashed"}
	db.Create(user)
	role := seedRegisterRole(db, "member", true)

	repo := register.NewRegisterRepository(db)
	err := repo.CreateUserRole(&models.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	})

	assert.NoError(t, err)
}

func TestIntegration_Register_CreateUserRole_Duplicate(t *testing.T) {
	db := setupRegisterTestDB(t)

	user := &models.User{Email: "newuser@example.com", Password: "hashed"}
	db.Create(user)
	role := seedRegisterRole(db, "member", true)

	repo := register.NewRegisterRepository(db)

	// Insert pertama
	_ = repo.CreateUserRole(&models.UserRole{UserID: user.ID, RoleID: role.ID})

	// Insert duplikat — harus error karena unique index
	err := repo.CreateUserRole(&models.UserRole{UserID: user.ID, RoleID: role.ID})

	assert.Error(t, err)
}

func TestIntegration_Register_FullFlow(t *testing.T) {
	db := setupRegisterTestDB(t)
	seedRegisterRole(db, "member", true)

	repo := register.NewRegisterRepository(db)

	// 1. Pastikan email belum ada
	_, err := repo.FindByEmail("newuser@example.com")
	assert.Error(t, err)

	// 2. Buat user baru
	user := &models.User{Email: "newuser@example.com", Password: "hashedpassword"}
	err = repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	// 3. Ambil role default
	role, err := repo.FindDefaultRole()
	assert.NoError(t, err)
	assert.NotZero(t, role.ID)

	// 4. Assign role
	err = repo.CreateUserRole(&models.UserRole{UserID: user.ID, RoleID: role.ID})
	assert.NoError(t, err)

	// 5. Verifikasi user tersimpan
	saved, err := repo.FindByEmail("newuser@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "newuser@example.com", saved.Email)
}
