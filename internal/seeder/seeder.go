package seeder

import (
	"log"

	"gorm.io/gorm"
)

var seeders = []func(*gorm.DB) error{
	SeedRoles,
	SeedPermissions,
	SeedUsers,
}

func Run(db *gorm.DB) error {
	log.Println("Running seeders...")
	for _, seed := range seeders {
		if err := seed(db); err != nil {
			return err
		}
	}
	log.Println("Seeders completed.")
	return nil
}
