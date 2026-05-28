package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name           string       `gorm:"unique;type:varchar(100);not null"`
	SystemFunction string       `gorm:"type:varchar(100)"`
	IsDefault      bool         `gorm:"default:false"`
	Permissions    []Permission `gorm:"foreignKey:RoleID"`
	Users          []UserRole   `gorm:"foreignKey:RoleID"`
}
