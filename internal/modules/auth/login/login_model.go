package login

import "gorm.io/gorm"


type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email    string `gorm:"unique;type:varchar(100);not null"`
	Password string `gorm:"type:varchar(255);not null"`
}
