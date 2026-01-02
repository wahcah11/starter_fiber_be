package profil

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
    loginRepo := login.NewLoginRepository(db)
    svc := NewProfilService(loginRepo)
    ctrl := NewProfilController(svc)

    auth := router.Group("/auth")
    auth.Get("/profil", middleware.Protected(), ctrl.GetProfile)
}
