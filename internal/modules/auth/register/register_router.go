package register

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewRegisterRepository(db)
	svc := NewRegisterService(repo)
	ctrl := NewRegisterController(svc)

	auth := router.Group("/auth")
	auth.Post("/register", ctrl.Register)
}
