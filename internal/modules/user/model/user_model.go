package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstName" gorm:"type:varchar(50);not null"`
	LastName  string `json:"lastName" gorm:"type:varchar(50);not null"`
	Role      string `json:"role" gorm:"type:varchar(20);default:'user'"`
	Email     string `json:"email" gorm:"unique;type:varchar(100);not null"`
	Password  string `json:"password" gorm:"type:varchar(255);not null"`
}
