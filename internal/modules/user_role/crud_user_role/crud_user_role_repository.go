package crud_user_role

import (
	"starter-wahcah-be/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll(page, limit int, userID *uint) ([]models.UserRole, int64, error)
	FindByID(id uint) (*models.UserRole, error)
	FindByUserIDAndRoleID(userID, roleID uint) (*models.UserRole, error)
	Create(userRole *models.UserRole) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewCrudUserRoleRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// FindAll mengambil semua user_role dengan pagination dan filter opsional by user_id
func (r *repository) FindAll(page, limit int, userID *uint) ([]models.UserRole, int64, error) {
	var userRoles []models.UserRole
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&models.UserRole{})

	// Filter by user_id jika diberikan
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// Hitung total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	if err := query.
		Offset(offset).
		Limit(limit).
		Order("user_id ASC, role_id ASC").
		Find(&userRoles).Error; err != nil {
		return nil, 0, err
	}

	return userRoles, total, nil
}

// FindByID mencari user_role berdasarkan ID
func (r *repository) FindByID(id uint) (*models.UserRole, error) {
	var userRole models.UserRole
	if err := r.db.First(&userRole, id).Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}

// FindByUserIDAndRoleID mencari user_role berdasarkan user_id dan role_id
// Digunakan untuk cek duplikat assignment
func (r *repository) FindByUserIDAndRoleID(userID, roleID uint) (*models.UserRole, error) {
	var userRole models.UserRole
	if err := r.db.
		Where("user_id = ? AND role_id = ?", userID, roleID).
		First(&userRole).Error; err != nil {
		return nil, err
	}
	return &userRole, nil
}

// Create menyimpan assignment role ke user
func (r *repository) Create(userRole *models.UserRole) error {
	return r.db.Create(userRole).Error
}

// Delete menghapus assignment role dari user berdasarkan ID (soft delete via gorm.Model)
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&models.UserRole{}, id).Error
}
