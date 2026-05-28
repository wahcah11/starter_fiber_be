package seeder

import (
	"log"
	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/util"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	users := []struct {
		Email    string
		Password string
		RoleName string
	}{
		{Email: "superadmin@example.com", Password: "superadmin123", RoleName: "superadmin"},
		{Email: "admin@example.com", Password: "admin123", RoleName: "admin"},
		{Email: "member@example.com", Password: "member123", RoleName: "member"},
	}

	for _, u := range users {
		var existing models.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			continue // Sudah ada, skip
		}

		hashedPassword, err := util.HashPassword(u.Password)
		if err != nil {
			return err
		}

		user := models.User{
			Email:    u.Email,
			Password: hashedPassword,
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}

		log.Printf("User created: %s", u.Email)

		// Assign role
		var role models.Role
		if err := db.Where("name = ?", u.RoleName).First(&role).Error; err != nil {
			log.Printf("Role '%s' not found, skip assign for %s", u.RoleName, u.Email)
			continue
		}

		if err := db.Create(&models.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}).Error; err != nil {
			return err
		}

		log.Printf("Role '%s' assigned to %s", u.RoleName, u.Email)
	}
	return nil
}
