package profile

import (
	"starter-wahcah-be/internal/modules/auth/login"

	"gorm.io/gorm"
)

type Repository interface {
	FindByID(userID uint) (*login.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByID(userID uint) (*login.User, error) {
	var user login.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
