package router

import (
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/profile"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	api := app.Group("/api")

	// Auth
	login.InitRoutes(api, db)
	// Profile
	profile.InitRoutes(api, db)
}
