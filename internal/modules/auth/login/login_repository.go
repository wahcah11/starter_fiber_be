package login

import "gorm.io/gorm"

type Repository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error) // ✅ TAMBAHKAN
	CreateUser(user *User) error // Tambahan untuk seeding/register manual
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
func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return &user, err
}
func (r *repository) CreateUser(user *User) error {
	// Simpan user baru ke database
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
