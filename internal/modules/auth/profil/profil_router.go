package profil

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"starter-wahcah-be/internal/middleware"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewProfilRepository(db)
	svc := NewProfilService(repo)
	ctrl := NewProfilController(svc)

	auth := router.Group("/auth")
	auth.Get("/profil", middleware.Protected(), ctrl.GetProfile)
}
