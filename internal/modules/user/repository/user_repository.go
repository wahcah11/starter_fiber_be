package repository

import (
	"errors"
	"starter-wahcah-be/internal/modules/user/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)

	if result.RowsAffected == 0 {
		return errors.New("User Not Found")
	}

	return result.Error
}
