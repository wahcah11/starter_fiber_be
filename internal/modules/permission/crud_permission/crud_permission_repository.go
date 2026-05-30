package crud_permission

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(page, limit int, roleID *uint) ([]models.Permission, int64, error)
	FindByID(id uint) (*models.Permission, error)
	FindByRoleIDAndName(roleID uint, name string) (*models.Permission, error)
	Create(permission *models.Permission) error
	Update(permission *models.Permission) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewCrudPermissionRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// FindAll mengambil semua permission dengan pagination dan filter opsional by role_id
func (r *repository) FindAll(page, limit int, roleID *uint) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.Permission{})

	// Filter by role_id jika diberikan
	if roleID != nil {
		query = query.Where("role_id = ?", *roleID)
	}

	// Hitung total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	if err := query.
		Offset(offset).
		Limit(limit).
		Order("role_id ASC, name ASC").
		Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// FindByID mencari permission berdasarkan ID
func (r *repository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByRoleIDAndName mencari permission berdasarkan role_id dan name
// Digunakan untuk cek duplikat
func (r *repository) FindByRoleIDAndName(roleID uint, name string) (*models.Permission, error) {
	var permission models.Permission
	if err := r.db.
		Where("role_id = ? AND name = ?", roleID, name).
		First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// Create menyimpan permission baru
func (r *repository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// Update menyimpan perubahan permission
func (r *repository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

// Delete menghapus permission berdasarkan ID (soft delete via gorm.Model)
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}
