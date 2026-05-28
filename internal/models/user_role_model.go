package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	UserID uint `gorm:"not null;uniqueIndex:idx_user_role"`
	RoleID uint `gorm:"not null;uniqueIndex:idx_user_role"`
	User   User `gorm:"foreignKey:UserID"`
	Role   Role `gorm:"foreignKey:RoleID"`
}
