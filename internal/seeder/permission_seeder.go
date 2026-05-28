package seeder

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) error {
	rolePermissions := map[string][]string{
		"superadmin": {
			"user:create", "user:read", "user:update", "user:delete",
			"role:create", "role:read", "role:update", "role:delete",
		},
		"admin": {
			"user:create", "user:read", "user:update",
			"role:read",
		},
		"member": {
			"user:read",
		},
	}

	for roleName, permissions := range rolePermissions {
		var role models.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			continue // Role belum ada, skip
		}

		for _, permName := range permissions {
			var existing models.Permission
			if err := db.Where("role_id = ? AND name = ?", role.ID, permName).First(&existing).Error; err != nil {
				db.Create(&models.Permission{RoleID: role.ID, Name: permName})
			}
		}
	}
	return nil
}
