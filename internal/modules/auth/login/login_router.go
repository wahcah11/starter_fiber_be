package login

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewLoginService(repo)
	ctrl := NewController(svc)

	auth := router.Group("/auth")
	auth.Post("/login", ctrl.Login)
	auth.Post("/register-test", ctrl.RegisterTest) // Endpoint sementara


}
