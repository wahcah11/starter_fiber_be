package router

import (
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	api := app.Group("/api")

	// Panggil Resepsionis Login
	login.InitRoutes(api, db)
}
