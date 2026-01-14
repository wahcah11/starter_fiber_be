package login

import "gorm.io/gorm"

type Repository interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (User, error)
	GetByID(id uint) (User, error)
	FindByEmail(email string) (*User, error)
	CreateUser(user *User) error // Tambahan untuk seeding/register manual
}

type repository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Register(req RegisterRequest) error {
	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password, // nanti bisa di-hash
	}

	return r.db.Create(&user).Error
}

// func (r *repository) Login(req LoginRequest) (User, error) {
//     var user User
//     err := r.db.Where("email = ?", req.Email).First(&user).Error
//     return user, err
// }

func (r *repository) Login(req LoginRequest) (User, error) {
	var user User
	err := r.db.Where("email = ?", req.Email).First(&user).Error
	//err := r.db.Where("email = ? AND password = ?", req.Email, req.Password).First(&user).Error
	return user, err
}

func (r *repository) GetByID(id uint) (User, error) {
	var user User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}
