package profile

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	// Inisialisasi repo, service, controller
	repo := login.NewLoginRepository(db)
	svc := NewProfileService(repo)
	ctrl := NewProfileController(svc)

	user := router.Group("/auth")
	user.Get("/getprofile", middleware.Protected(), ctrl.GetProfile)
}
