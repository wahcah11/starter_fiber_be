package router

import (
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/profile"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")

	//login
	login.InitRoutes(api, db)

	// profil
	profile.InitRoutes(api, db)
}