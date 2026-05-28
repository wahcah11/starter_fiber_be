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
	v1 := api.Group("/v1")

	// Panggil Resepsionis Login
	login.InitRoutes(v1, db)
}
