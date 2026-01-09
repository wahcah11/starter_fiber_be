package login

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	//repo := NewLoginRepository(db)
	repo := NewRepository(db)
	svc := NewLoginService(repo)
	ctrl := NewController(svc)
	//ctrl := NewLoginController(svc)

	auth := router.Group("/auth")
	auth.Post("/login", ctrl.Login)
	auth.Post("/register-test", ctrl.RegisterTest) // Endpoint sementara

	// profile := router.Group("/auth")
	// profile.Get("/profil", ctrl.GetProfile)
}
