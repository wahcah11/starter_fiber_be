package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	RoleID uint   `gorm:"not null;uniqueIndex:idx_role_permission"`
	Name   string `gorm:"type:varchar(100);not null;uniqueIndex:idx_role_permission"`
	Role   Role   `gorm:"foreignKey:RoleID"`
}
