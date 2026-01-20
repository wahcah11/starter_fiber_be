package profile

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(api fiber.Router, loginRepo login.Repository) {
	service := NewProfileService(loginRepo)
	controller := NewProfileController(service)

	profileRoute := api.Group(
		"/profile",
		middleware.Protected(),
	)

	profileRoute.Get("/", controller.GetProfile)
}
