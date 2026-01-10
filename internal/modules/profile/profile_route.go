package profile

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login"
	//"starter-wahcah-be/internal/modules/auth/profile"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {

	// Pakai repository login untuk mengambil data user
	loginRepo := login.NewRepository(db)

	// Pass loginRepo ke profil service
	profilService := NewProfileService(loginRepo)

	// Controller profil
	profilController := NewProfileController(profilService)

	// Route: /api/auth/profil
	r := router.Group("/auth", middleware.AuthMiddleware)

	r.Get("/profil", profilController.GetProfile)
}
