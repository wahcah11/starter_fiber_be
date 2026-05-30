package crud

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(page, limit int) ([]models.Role, int64, error)
	FindByID(id uint) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewCrudRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// FindAll mengambil semua role dengan pagination
func (r *repository) FindAll(page, limit int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	offset := (page - 1) * limit

	// Hitung total data
	if err := r.db.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	if err := r.db.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// FindByID mencari role berdasarkan ID
func (r *repository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName mencari role berdasarkan nama
func (r *repository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// Create menyimpan role baru
func (r *repository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// Update menyimpan perubahan role
func (r *repository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete menghapus role berdasarkan ID (soft delete via gorm.Model)
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
