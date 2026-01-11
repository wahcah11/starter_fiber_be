package profile

import (
    "errors"
    "strings"

    "starter-wahcah-be/internal/modules/auth/login"
)

type Service interface {
    GetProfile(userID uint) (*ProfileResponse, error)
}

type service struct {
    repo interface {
        FindByID(id uint) (*login.User, error)
    }
}

func NewProfileService(repo interface {
    FindByID(id uint) (*login.User, error)
}) Service {
    return &service{repo: repo}
}

func (s *service) GetProfile(userID uint) (*ProfileResponse, error) {
    user, err := s.repo.FindByID(userID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    // ✅ Ambil full_name dari kolom Nama
    fullName := user.Nama

    // ✅ Fallback kalau Nama kosong (opsional)
    if strings.TrimSpace(fullName) == "" {
        fullName = strings.TrimSpace(user.FirstName + " " + user.LastName)
    }

    return &ProfileResponse{
        Nama : fullName,
        Email:    user.Email,
    }, nil
}
