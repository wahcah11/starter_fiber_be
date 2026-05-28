package seeder

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "superadmin", SystemFunction: "full_access", IsDefault: false},
		{Name: "admin", SystemFunction: "admin_access", IsDefault: false},
		{Name: "member", SystemFunction: "basic_access", IsDefault: true},
	}

	for _, role := range roles {
		var existing models.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			// Belum ada, buat baru
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
		// Sudah ada, skip (idempotent)
	}
	return nil
}
