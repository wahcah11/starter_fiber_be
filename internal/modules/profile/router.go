package profile

import (
	"starter-wahcah-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	// Inisialisasi repo, service, controller
	repo := NewProfileRepository(db)
	svc := NewProfileService(repo)
	ctrl := NewProfileController(svc)

	user := router.Group("/user")
	user.Get("/profile", middleware.JWTMiddleware, ctrl.GetProfile)
}
