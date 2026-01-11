package login

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(50);not null"`
	LastName  string `gorm:"type:varchar(50);not null"`
	Email    string `gorm:"unique;type:varchar(100);not null"`
	Password string `gorm:"type:varchar(255);not null"`
}
