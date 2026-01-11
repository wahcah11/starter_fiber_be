package login

import (
	"errors"
	"starter-wahcah-be/internal/util"
)

type Service interface {
	Authenticate(req LoginRequest) (*LoginResponse, error)
	RegisterUser(req RegisterRequest) error // Menggunakan RegisterRequest
}

type service struct {
	repo Repository
}

func NewLoginService(repo Repository) Service {
	return &service{repo}
}

// Fungsi Authenticate untuk autentikasi pengguna
func (s *service) Authenticate(req LoginRequest) (*LoginResponse, error) {
	// Mencari pengguna berdasarkan email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Memeriksa apakah password yang dimasukkan benar
	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Membuat token untuk autentikasi
	token, _ := util.GenerateToken(user.ID)
	return &LoginResponse{Token: token}, nil
}

// Fungsi RegisterUser untuk mendaftar pengguna baru
// Fungsi RegisterUser untuk mendaftar pengguna baru
func (s *service) RegisterUser(req RegisterRequest) error {
	// Cek apakah email sudah terdaftar
	existingUser, err := s.repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		// Jika user sudah ada, kembalikan error
		return errors.New("email already in use")
	}

	// Meng-hash password yang diterima dari request
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		// Tangani error jika password gagal di-hash
		return errors.New("failed to hash password")
	}

	// Gabungkan first_name dan last_name menjadi nama lengkap
	namaLengkap := req.FirstName + " " + req.LastName

	// Membuat objek User baru berdasarkan data yang diterima
	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Nama:      namaLengkap, // Isi kolom Nama dengan nama lengkap
	}

	// Menyimpan user baru ke database dan menangani error
	err = s.repo.CreateUser(&user)
	if err != nil {
		// Tangani error jika gagal menyimpan user ke database
		return errors.New("failed to create user")
	}

	return nil
}
