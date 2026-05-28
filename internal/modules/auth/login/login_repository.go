package login

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
