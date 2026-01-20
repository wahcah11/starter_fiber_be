package router

import (
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/auth/profile"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	api := app.Group("/api")

	// ✅ gunakan nama constructor yang ada di login_repository.go
	loginRepo := login.NewLoginRepository(db)

	// ✅ inject repo yang sama
	login.InitRoutes(api, loginRepo)
	profile.InitRoutes(api, loginRepo)
}
