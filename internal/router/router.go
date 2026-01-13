package router

import (
	"starter-wahcah-be/internal/modules/auth/login"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	
    "starter-wahcah-be/internal/modules/auth/profile"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
    app.Use(logger.New())

    api := app.Group("/api")

    login.InitRoutes(api, db)
    profile.InitRoutes(api, db)
}

