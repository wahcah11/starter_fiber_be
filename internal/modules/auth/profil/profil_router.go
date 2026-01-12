package profil

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	loginRepo := login.NewLoginRepository(db)

	profilService := NewProfilService(loginRepo)
	profilController := NewProfilController(profilService)

	auth := router.Group("/auth")
	auth.Get("/profil", middleware.Protected(), profilController.GetProfile)
}
