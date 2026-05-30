package crud_test

import (
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/role/crud"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =========================================
// Helper
// =========================================

func setupCrudTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&models.Role{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func seedRole(db *gorm.DB, name, systemFunction string, isDefault bool) *models.Role {
	role := &models.Role{
		Name:           name,
		SystemFunction: systemFunction,
		IsDefault:      isDefault,
	}
	db.Create(role)
	return role
}

// =========================================
// FindAll
// =========================================

func TestIntegration_Crud_FindAll_Success(t *testing.T) {
	db := setupCrudTestDB(t)
	seedRole(db, "superadmin", "full_access", false)
	seedRole(db, "admin", "admin_access", false)
	seedRole(db, "member", "basic_access", true)

	repo := crud.NewCrudRepository(db)
	result, total, err := repo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 3, len(result))
}

func TestIntegration_Crud_FindAll_Pagination(t *testing.T) {
	db := setupCrudTestDB(t)
	seedRole(db, "role1", "function1", false)
	seedRole(db, "role2", "function2", false)
	seedRole(db, "role3", "function3", false)
	seedRole(db, "role4", "function4", false)
	seedRole(db, "role5", "function5", false)

	repo := crud.NewCrudRepository(db)

	// Halaman 1, limit 2 — harus dapat 2 data
	result, total, err := repo.FindAll(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Equal(t, 2, len(result))

	// Halaman 2, limit 2 — harus dapat 2 data
	result2, _, err := repo.FindAll(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result2))

	// Halaman 3, limit 2 — harus dapat 1 data (sisa)
	result3, _, err := repo.FindAll(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result3))
}

func TestIntegration_Crud_FindAll_Empty(t *testing.T) {
	db := setupCrudTestDB(t)

	repo := crud.NewCrudRepository(db)
	result, total, err := repo.FindAll(1, 10)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(result))
}

// =========================================
// FindByID
// =========================================

func TestIntegration_Crud_FindByID_Success(t *testing.T) {
	db := setupCrudTestDB(t)
	role := seedRole(db, "superadmin", "full_access", false)

	repo := crud.NewCrudRepository(db)
	result, err := repo.FindByID(role.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "superadmin", result.Name)
}

func TestIntegration_Crud_FindByID_NotFound(t *testing.T) {
	db := setupCrudTestDB(t)

	repo := crud.NewCrudRepository(db)
	result, err := repo.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// =========================================
// FindByName
// =========================================

func TestIntegration_Crud_FindByName_Success(t *testing.T) {
	db := setupCrudTestDB(t)
	seedRole(db, "member", "basic_access", true)

	repo := crud.NewCrudRepository(db)
	result, err := repo.FindByName("member")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "member", result.Name)
	assert.True(t, result.IsDefault)
}

func TestIntegration_Crud_FindByName_NotFound(t *testing.T) {
	db := setupCrudTestDB(t)

	repo := crud.NewCrudRepository(db)
	result, err := repo.FindByName("unknown")

	assert.Error(t, err)
	assert.Nil(t, result)
}

// =========================================
// Create
// =========================================

func TestIntegration_Crud_Create_Success(t *testing.T) {
	db := setupCrudTestDB(t)

	repo := crud.NewCrudRepository(db)
	role := &models.Role{
		Name:           "editor",
		SystemFunction: "edit_access",
		IsDefault:      false,
	}

	err := repo.Create(role)

	assert.NoError(t, err)
	assert.NotZero(t, role.ID)
}

func TestIntegration_Crud_Create_DuplicateName(t *testing.T) {
	db := setupCrudTestDB(t)
	seedRole(db, "editor", "edit_access", false)

	repo := crud.NewCrudRepository(db)
	duplicate := &models.Role{
		Name:           "editor",
		SystemFunction: "edit_access",
	}

	err := repo.Create(duplicate)

	// Harus error karena name unique
	assert.Error(t, err)
}

// =========================================
// Update
// =========================================

func TestIntegration_Crud_Update_Success(t *testing.T) {
	db := setupCrudTestDB(t)
	role := seedRole(db, "editor", "edit_access", false)

	repo := crud.NewCrudRepository(db)
	role.Name = "editor-updated"
	role.SystemFunction = "edit_access_v2"

	err := repo.Update(role)

	assert.NoError(t, err)

	// Verifikasi perubahan tersimpan
	updated, _ := repo.FindByID(role.ID)
	assert.Equal(t, "editor-updated", updated.Name)
	assert.Equal(t, "edit_access_v2", updated.SystemFunction)
}

// =========================================
// Delete
// =========================================

func TestIntegration_Crud_Delete_Success(t *testing.T) {
	db := setupCrudTestDB(t)
	role := seedRole(db, "editor", "edit_access", false)

	repo := crud.NewCrudRepository(db)
	err := repo.Delete(role.ID)

	assert.NoError(t, err)

	// Verifikasi sudah terhapus (soft delete)
	_, err = repo.FindByID(role.ID)
	assert.Error(t, err)
}

func TestIntegration_Crud_Delete_NotFound(t *testing.T) {
	db := setupCrudTestDB(t)

	repo := crud.NewCrudRepository(db)
	// GORM soft delete tidak return error jika ID tidak ditemukan
	// Test ini sebagai dokumentasi behavior
	err := repo.Delete(999)
	assert.NoError(t, err)
}

// =========================================
// Full Flow
// =========================================

func TestIntegration_Crud_FullFlow(t *testing.T) {
	db := setupCrudTestDB(t)
	repo := crud.NewCrudRepository(db)

	// 1. Create
	role := &models.Role{Name: "editor", SystemFunction: "edit_access", IsDefault: false}
	err := repo.Create(role)
	assert.NoError(t, err)
	assert.NotZero(t, role.ID)

	// 2. FindByID
	found, err := repo.FindByID(role.ID)
	assert.NoError(t, err)
	assert.Equal(t, "editor", found.Name)

	// 3. FindByName
	byName, err := repo.FindByName("editor")
	assert.NoError(t, err)
	assert.Equal(t, role.ID, byName.ID)

	// 4. Update
	role.Name = "editor-updated"
	err = repo.Update(role)
	assert.NoError(t, err)

	updated, _ := repo.FindByID(role.ID)
	assert.Equal(t, "editor-updated", updated.Name)

	// 5. FindAll
	roles, total, err := repo.FindAll(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(roles))

	// 6. Delete
	err = repo.Delete(role.ID)
	assert.NoError(t, err)

	// 7. Verifikasi terhapus
	_, err = repo.FindByID(role.ID)
	assert.Error(t, err)
}
