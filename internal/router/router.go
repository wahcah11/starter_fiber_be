package router

import (
	"starter-wahcah-be/internal/modules/auth/login"
	//"starter-wahcah-be/internal/modules/profil"
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	app.Use(logger.New())

	api := app.Group("/api")

	// LOGIN MODULE
	login.InitRoutes(api, db)

	// PROFILE MODULE
	profileRepo := profile.NewRepository(db)
	profileService := profile.NewService(profileRepo)
	profileController := profile.NewController(profileService)

	// ROUTE PROFIL
	profile := api.Group("/profil", middleware.AuthMiddleware)
	profile.Get("/me", profileController.GetProfile)
}

// func SetupRoutes(app *fiber.App, db *gorm.DB) {
// 	repo := login.NewRepository(db)
// 	svc := login.NewLoginService(repo)
// 	controller := login.NewController(svc)

// 	api := app.Group("/api/auth")

// 	api.Post("/register-test", controller.Register)
// 	api.Post("/login", controller.Login)
// 	api.Get("/profil", middleware.AuthMiddleware, controller.Profil)
// }

// func SetupRoutes(app *fiber.App, db *gorm.DB) {
// 	app.Use(logger.New())

// 	api := app.Group("/api")

// 	// Panggil Resepsionis Login
// 	login.InitRoutes(api, db)
// }


// func SetupRoutes(app *fiber.App, db *gorm.DB) {

// 	app.Use(logger.New())
// 	// PROFILE MODULE
// 	profileRepo := profile.NewRepository(db)
// 	profileService := profile.NewService(profileRepo)
// 	profileController := profile.NewController(profileService)

// 	api := app.Group("/api")

// 	login.InitRoutes(api, db)
// 	auth := api.Group("/auth")
// 	auth.Get("/profil", middleware.AuthMiddleware, profileController.GetProfile)
// }

// func SetupRoutes(app *fiber.App, db *gorm.DB) {
// 	app.Use(logger.New())

// 	api := app.Group("/api")

// 	// Panggil Resepsionis Login
// 	login.InitRoutes(api, db)
// }
