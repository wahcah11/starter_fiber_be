package models

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Role{},
		&Permission{},
		&User{},
		&UserRole{},
	)
}
