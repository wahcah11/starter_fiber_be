package login

import "gorm.io/gorm"

type User struct {
	gorm.Model
FirstName string `json:"first_name" gorm:"type:varchar(100);not null"`
LastName  string `json:"last_name" gorm:"type:varchar(100);not null"`
Nama      string `json:"nama" gorm:"type:varchar(255);not null"`
Email     string `json:"email" gorm:"unique;type:varchar(100);not null"`
Password  string `json:"password" gorm:"type:varchar(255);not null"`

}
