package profile

import "gorm.io/gorm"

type Repository interface {
	GetByID(id uint) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetByID(id uint) (User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return user, err
}
