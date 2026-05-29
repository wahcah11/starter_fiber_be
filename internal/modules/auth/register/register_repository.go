package register

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	FindDefaultRole() (*models.Role, error)
	CreateUserRole(userRole *models.UserRole) error
}

type repository struct {
	db *gorm.DB
}

func NewRegisterRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// FindByEmail mencari user berdasarkan email
func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Create menyimpan user baru ke database
func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindDefaultRole mengambil role dengan is_default = true
func (r *repository) FindDefaultRole() (*models.Role, error) {
	var role models.Role
	err := r.db.Where("is_default = ?", true).First(&role).Error
	return &role, err
}

// CreateUserRole menyimpan relasi user dan role
func (r *repository) CreateUserRole(userRole *models.UserRole) error {
	return r.db.Create(userRole).Error
}
