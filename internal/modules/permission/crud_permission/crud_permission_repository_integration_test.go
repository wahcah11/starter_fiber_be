package crud_permission_test

import (
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/permission/crud_permission"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =========================================
// Helper
// =========================================

func setupCrudPermissionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
	); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func seedPermissionRole(db *gorm.DB, name string) *models.Role {
	role := &models.Role{Name: name, SystemFunction: "test_function"}
	db.Create(role)
	return role
}

func seedPermission(db *gorm.DB, roleID uint, name string) *models.Permission {
	permission := &models.Permission{RoleID: roleID, Name: name}
	db.Create(permission)
	return permission
}

// =========================================
// FindAll
// =========================================

func TestIntegration_Permission_FindAll_Success(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role1 := seedPermissionRole(db, "superadmin")
	role2 := seedPermissionRole(db, "admin")
	seedPermission(db, role1.ID, "role:read")
	seedPermission(db, role1.ID, "role:create")
	seedPermission(db, role2.ID, "user:read")

	repo := crud_permission.NewCrudPermissionRepository(db)
	result, total, err := repo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 3, len(result))
}

func TestIntegration_Permission_FindAll_FilterByRoleID(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role1 := seedPermissionRole(db, "superadmin")
	role2 := seedPermissionRole(db, "admin")
	seedPermission(db, role1.ID, "role:read")
	seedPermission(db, role1.ID, "role:create")
	seedPermission(db, role2.ID, "user:read")

	repo := crud_permission.NewCrudPermissionRepository(db)

	// Filter hanya permission milik role1
	result, total, err := repo.FindAll(1, 10, &role1.ID)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(result))
	for _, p := range result {
		assert.Equal(t, role1.ID, p.RoleID)
	}
}

func TestIntegration_Permission_FindAll_Pagination(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	seedPermission(db, role.ID, "perm:1")
	seedPermission(db, role.ID, "perm:2")
	seedPermission(db, role.ID, "perm:3")
	seedPermission(db, role.ID, "perm:4")
	seedPermission(db, role.ID, "perm:5")

	repo := crud_permission.NewCrudPermissionRepository(db)

	// Halaman 1, limit 2
	result1, total, err := repo.FindAll(1, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)
	assert.Equal(t, 2, len(result1))

	// Halaman 2, limit 2
	result2, _, err := repo.FindAll(2, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result2))

	// Halaman 3, limit 2 — sisa 1
	result3, _, err := repo.FindAll(3, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result3))
}

func TestIntegration_Permission_FindAll_Empty(t *testing.T) {
	db := setupCrudPermissionTestDB(t)

	repo := crud_permission.NewCrudPermissionRepository(db)
	result, total, err := repo.FindAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Equal(t, 0, len(result))
}

// =========================================
// FindByID
// =========================================

func TestIntegration_Permission_FindByID_Success(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	permission := seedPermission(db, role.ID, "role:read")

	repo := crud_permission.NewCrudPermissionRepository(db)
	result, err := repo.FindByID(permission.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "role:read", result.Name)
	assert.Equal(t, role.ID, result.RoleID)
}

func TestIntegration_Permission_FindByID_NotFound(t *testing.T) {
	db := setupCrudPermissionTestDB(t)

	repo := crud_permission.NewCrudPermissionRepository(db)
	result, err := repo.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// =========================================
// FindByRoleIDAndName
// =========================================

func TestIntegration_Permission_FindByRoleIDAndName_Success(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	seedPermission(db, role.ID, "role:read")

	repo := crud_permission.NewCrudPermissionRepository(db)
	result, err := repo.FindByRoleIDAndName(role.ID, "role:read")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "role:read", result.Name)
}

func TestIntegration_Permission_FindByRoleIDAndName_NotFound(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")

	repo := crud_permission.NewCrudPermissionRepository(db)
	result, err := repo.FindByRoleIDAndName(role.ID, "role:delete")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_Permission_FindByRoleIDAndName_SameNameDifferentRole(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role1 := seedPermissionRole(db, "superadmin")
	role2 := seedPermissionRole(db, "admin")

	// Nama sama tapi role berbeda
	seedPermission(db, role1.ID, "role:read")
	seedPermission(db, role2.ID, "role:read")

	repo := crud_permission.NewCrudPermissionRepository(db)

	// Harus return permission milik role1 saja
	result, err := repo.FindByRoleIDAndName(role1.ID, "role:read")
	assert.NoError(t, err)
	assert.Equal(t, role1.ID, result.RoleID)

	// Harus return permission milik role2 saja
	result2, err := repo.FindByRoleIDAndName(role2.ID, "role:read")
	assert.NoError(t, err)
	assert.Equal(t, role2.ID, result2.RoleID)
}

// =========================================
// Create
// =========================================

func TestIntegration_Permission_Create_Success(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")

	repo := crud_permission.NewCrudPermissionRepository(db)
	permission := &models.Permission{RoleID: role.ID, Name: "role:delete"}

	err := repo.Create(permission)

	assert.NoError(t, err)
	assert.NotZero(t, permission.ID)
}

func TestIntegration_Permission_Create_Duplicate(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	seedPermission(db, role.ID, "role:read")

	repo := crud_permission.NewCrudPermissionRepository(db)
	duplicate := &models.Permission{RoleID: role.ID, Name: "role:read"}

	// Harus error karena unique index (role_id, name)
	err := repo.Create(duplicate)
	assert.Error(t, err)
}

// =========================================
// Update
// =========================================

func TestIntegration_Permission_Update_Success(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	permission := seedPermission(db, role.ID, "role:read")

	repo := crud_permission.NewCrudPermissionRepository(db)
	permission.Name = "role:read-updated"

	err := repo.Update(permission)

	assert.NoError(t, err)

	// Verifikasi perubahan tersimpan
	updated, _ := repo.FindByID(permission.ID)
	assert.Equal(t, "role:read-updated", updated.Name)
}

// =========================================
// Delete
// =========================================

func TestIntegration_Permission_Delete_Success(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	permission := seedPermission(db, role.ID, "role:read")

	repo := crud_permission.NewCrudPermissionRepository(db)
	err := repo.Delete(permission.ID)

	assert.NoError(t, err)

	// Verifikasi sudah terhapus
	_, err = repo.FindByID(permission.ID)
	assert.Error(t, err)
}

func TestIntegration_Permission_Delete_NotFound(t *testing.T) {
	db := setupCrudPermissionTestDB(t)

	repo := crud_permission.NewCrudPermissionRepository(db)
	// GORM soft delete tidak return error jika ID tidak ada
	err := repo.Delete(999)
	assert.NoError(t, err)
}

// =========================================
// Full Flow
// =========================================

func TestIntegration_Permission_FullFlow(t *testing.T) {
	db := setupCrudPermissionTestDB(t)
	role := seedPermissionRole(db, "superadmin")
	repo := crud_permission.NewCrudPermissionRepository(db)

	// 1. Create
	permission := &models.Permission{RoleID: role.ID, Name: "role:read"}
	err := repo.Create(permission)
	assert.NoError(t, err)
	assert.NotZero(t, permission.ID)

	// 2. FindByID
	found, err := repo.FindByID(permission.ID)
	assert.NoError(t, err)
	assert.Equal(t, "role:read", found.Name)

	// 3. FindByRoleIDAndName
	byName, err := repo.FindByRoleIDAndName(role.ID, "role:read")
	assert.NoError(t, err)
	assert.Equal(t, permission.ID, byName.ID)

	// 4. Update
	permission.Name = "role:read-updated"
	err = repo.Update(permission)
	assert.NoError(t, err)

	updated, _ := repo.FindByID(permission.ID)
	assert.Equal(t, "role:read-updated", updated.Name)

	// 5. FindAll
	perms, total, err := repo.FindAll(1, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, 1, len(perms))

	// 6. FindAll filter by role
	permsFiltered, totalFiltered, err := repo.FindAll(1, 10, &role.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalFiltered)
	assert.Equal(t, 1, len(permsFiltered))

	// 7. Delete
	err = repo.Delete(permission.ID)
	assert.NoError(t, err)

	// 8. Verifikasi terhapus
	_, err = repo.FindByID(permission.ID)
	assert.Error(t, err)
}
