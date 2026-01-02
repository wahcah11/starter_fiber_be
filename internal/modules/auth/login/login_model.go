package login

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Firstname string `gorm:"type:varchar(100);not null"`
    Lastname  string `gorm:"type:varchar(100);not null"`
    Email     string `gorm:"unique;type:varchar(100);not null"`
    Password  string `gorm:"type:varchar(255);not null"`
}
