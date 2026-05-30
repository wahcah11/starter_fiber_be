package crud_user_role_test

import (
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/user_role/crud_user_role"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =========================================
// Helper
// =========================================

func setupCrudUserRoleTestDB(t *testing.T) *gorm.DB {
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

func seedUserRoleUser(db *gorm.DB, email string) *models.User {
	user := &models.User{Email: email, Password: "hashed"}
	db.Create(user)
	return user
}

func seedUserRoleRole(db *gorm.DB, name string) *models.Role {
	role := &models.Role{Name: name, SystemFunction: "test_function"}
	db.Create(role)
	return role
}

func seedUserRole(db *gorm.DB, userID, roleID uint) *models.UserRole {
	userRole := &models.UserRole{UserID: userID, RoleID: roleID}
	db.Create(userRole)
	return userRole
}

// =========================================
// FindAll
// =========================================

func TestIntegration_UserRole_FindAll_Success(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user1 := seedUserRoleUser(db, "user1@example.com")
	user2 := seedUserRoleUser(db, "user2@example.com")
	role1 := seedUserRoleRole(db, "admin")
	role2 := seedUserRoleRole(db, "member")
	seedUserRole(db, user1.ID, role1.ID)
	seedUserRole(db, user1.ID, role2.ID)
	seedUserRole(db, user2.ID, role1.ID)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	result, total, err := repo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 3, len(result))
}

func TestIntegration_UserRole_FindAll_FilterByUserID(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user1 := seedUserRoleUser(db, "user1@example.com")
	user2 := seedUserRoleUser(db, "user2@example.com")
	role1 := seedUserRoleRole(db, "admin")
	role2 := seedUserRoleRole(db, "member")
	seedUserRole(db, user1.ID, role1.ID)
	seedUserRole(db, user1.ID, role2.ID)
	seedUserRole(db, user2.ID, role1.ID)

	repo := crud_user_role.NewCrudUserRoleRepository(db)

	// Filter hanya user_role milik user1
	result, total, err := repo.FindAll(1, 10, &user1.ID)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))
	for _, ur := range result {
		assert.Equal(t, user1.ID, ur.UserID)
	}
}

func TestIntegration_UserRole_FindAll_Pagination(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	for i := 1; i <= 5; i++ {
		role := seedUserRoleRole(db, "role"+string(rune('0'+i)))
		seedUserRole(db, user.ID, role.ID)
	}

	repo := crud_user_role.NewCrudUserRoleRepository(db)

	// Halaman 1 limit 2
	result1, total, err := repo.FindAll(1, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Equal(t, 2, len(result1))

	// Halaman 2 limit 2
	result2, _, err := repo.FindAll(2, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result2))

	// Halaman 3 limit 2 — sisa 1
	result3, _, err := repo.FindAll(3, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result3))
}

func TestIntegration_UserRole_FindAll_Empty(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	result, total, err := repo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(result))
}

// =========================================
// FindByID
// =========================================

func TestIntegration_UserRole_FindByID_Success(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	role := seedUserRoleRole(db, "admin")
	userRole := seedUserRole(db, user.ID, role.ID)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	result, err := repo.FindByID(userRole.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, role.ID, result.RoleID)
}

func TestIntegration_UserRole_FindByID_NotFound(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	result, err := repo.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// =========================================
// FindByUserIDAndRoleID
// =========================================

func TestIntegration_UserRole_FindByUserIDAndRoleID_Success(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	role := seedUserRoleRole(db, "admin")
	seedUserRole(db, user.ID, role.ID)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	result, err := repo.FindByUserIDAndRoleID(user.ID, role.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.UserID)
	assert.Equal(t, role.ID, result.RoleID)
}

func TestIntegration_UserRole_FindByUserIDAndRoleID_NotFound(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	result, err := repo.FindByUserIDAndRoleID(user.ID, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// =========================================
// Create
// =========================================

func TestIntegration_UserRole_Create_Success(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	role := seedUserRoleRole(db, "admin")

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	userRole := &models.UserRole{UserID: user.ID, RoleID: role.ID}

	err := repo.Create(userRole)

	assert.NoError(t, err)
	assert.NotZero(t, userRole.ID)
}

func TestIntegration_UserRole_Create_Duplicate(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	role := seedUserRoleRole(db, "admin")
	seedUserRole(db, user.ID, role.ID)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	duplicate := &models.UserRole{UserID: user.ID, RoleID: role.ID}

	// Harus error karena unique index (user_id, role_id)
	err := repo.Create(duplicate)
	assert.Error(t, err)
}

// =========================================
// Delete
// =========================================

func TestIntegration_UserRole_Delete_Success(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	role := seedUserRoleRole(db, "admin")
	userRole := seedUserRole(db, user.ID, role.ID)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	err := repo.Delete(userRole.ID)

	assert.NoError(t, err)

	// Verifikasi sudah terhapus
	_, err = repo.FindByID(userRole.ID)
	assert.Error(t, err)
}

func TestIntegration_UserRole_Delete_NotFound(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)

	repo := crud_user_role.NewCrudUserRoleRepository(db)
	// GORM soft delete tidak return error jika ID tidak ada
	err := repo.Delete(999)
	assert.NoError(t, err)
}

// =========================================
// Full Flow
// =========================================

func TestIntegration_UserRole_FullFlow(t *testing.T) {
	db := setupCrudUserRoleTestDB(t)
	user := seedUserRoleUser(db, "user@example.com")
	role := seedUserRoleRole(db, "admin")
	repo := crud_user_role.NewCrudUserRoleRepository(db)

	// 1. Pastikan belum ada
	_, err := repo.FindByUserIDAndRoleID(user.ID, role.ID)
	assert.Error(t, err)

	// 2. Create
	userRole := &models.UserRole{UserID: user.ID, RoleID: role.ID}
	err = repo.Create(userRole)
	assert.NoError(t, err)
	assert.NotZero(t, userRole.ID)

	// 3. FindByID
	found, err := repo.FindByID(userRole.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, found.UserID)
	assert.Equal(t, role.ID, found.RoleID)

	// 4. FindByUserIDAndRoleID
	byUserRole, err := repo.FindByUserIDAndRoleID(user.ID, role.ID)
	assert.NoError(t, err)
	assert.Equal(t, userRole.ID, byUserRole.ID)

	// 5. FindAll
	userRoles, total, err := repo.FindAll(1, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(userRoles))

	// 6. FindAll filter by user_id
	filtered, filteredTotal, err := repo.FindAll(1, 10, &user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), filteredTotal)
	assert.Equal(t, 1, len(filtered))

	// 7. Delete
	err = repo.Delete(userRole.ID)
	assert.NoError(t, err)

	// 8. Verifikasi terhapus
	_, err = repo.FindByID(userRole.ID)
	assert.Error(t, err)
}
