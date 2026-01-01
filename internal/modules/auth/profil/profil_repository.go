package profil

import "gorm.io/gorm"

type Repository interface {
	FindByID(userID uint) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewProfilRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindByID(userID uint) (*User, error) {
	var user User
	err := r.db.First(&user, userID).Error
	return &user, err
}
