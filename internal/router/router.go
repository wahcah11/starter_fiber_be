package router

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/profile"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")

	//login
	login.InitRoutes(api, db)

	// profil
	profile.InitRoutes(api, db)
}
