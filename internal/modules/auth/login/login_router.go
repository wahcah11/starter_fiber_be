package login

import "github.com/gofiber/fiber/v2"

func InitRoutes(router fiber.Router, repo Repository) {
	svc := NewLoginService(repo)
	ctrl := NewLoginController(svc)

	auth := router.Group("/auth")
	auth.Post("/login", ctrl.Login)
	auth.Post("/register-test", ctrl.RegisterTest) // endpoint sementara
}
