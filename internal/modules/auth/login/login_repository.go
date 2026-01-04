package login

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*User, error)
	CreateUser(user *User) error // Tambahan untuk seeding/register manual
	FindByID(userID uint) (*User, error)

}

type repository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}


// For Profile Repository
func (r *repository) FindByID(userID uint) (*User, error) {
	var user User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
