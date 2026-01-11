package login

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewLoginRepository(db)
	svc := NewLoginService(repo)
	ctrl := NewLoginController(svc)

	auth := router.Group("/auth")
	auth.Post("/login", ctrl.Login)
	auth.Post("/register-test", ctrl.RegisterTest)
	 // Endpoint sementara
}
