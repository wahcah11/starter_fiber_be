package profile

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/auth/login" // ✅ WAJIB

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func InitRoutes(api fiber.Router, db *gorm.DB) {

	// 🔁 pakai repository login
	loginRepo := login.NewLoginRepository(db)

	// service profile pakai repo login
	service := NewProfileService(loginRepo)
	controller := NewProfileController(service)

	profile := api.Group(
		"/profile",
		middleware.Protected(),
	)

	profile.Get("/", controller.GetProfile)
}
