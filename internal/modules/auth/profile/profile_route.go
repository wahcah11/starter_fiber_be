package profile

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {

	loginRepo := login.NewLoginRepository(db)
	profilService := NewProfileService(loginRepo)
	profilController := NewProfileController(profilService)


	r := router.Group("/auth", middleware.Protected())
	r.Get("/profile", profilController.GetProfile)
}
