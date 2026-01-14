package profile

import (
	"starter-wahcah-be/internal/modules/auth/login"

	"gorm.io/gorm"
)

type Repository interface {
	GetByID(id uint) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByID(id uint) (User, error) {
	var user User
	var userData login.User

	err := r.db.First(&userData, id).Error
	if err != nil {
		return user, err
	}

	user.ID = userData.ID
	user.FirstName = userData.FirstName
	user.LastName = userData.LastName
	user.Email = userData.Email

	return user, nil
}
