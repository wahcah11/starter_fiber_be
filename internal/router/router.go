package router

import (
	"starter-wahcah-be/internal/modules/auth/login"
	//"starter-wahcah-be/internal/modules/profil"
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/profile"
	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

// func SetupRoutes(app *fiber.App, db *gorm.DB) {

// 	app.Use(logger.New())

// 	api := app.Group("/api")

// 	// LOGIN MODULE
// 	login.InitRoutes(api, db)

// 	// PROFILE MODULE
// 	profileRepo := profile.NewRepository(db)
// 	profileService := profile.NewService(profileRepo)
// 	profileController := profile.NewController(profileService)

// 	// ROUTE PROFIL
// 	profile := api.Group("/profil", middleware.AuthMiddleware)
// 	profile.Get("/me", profileController.GetProfile)
// }

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// --- LOGIN MODULE ---
	loginRepo := login.NewRepository(db)
	loginSvc := login.NewLoginService(loginRepo)
	loginCtrl := login.NewController(loginSvc)

	// --- PROFILE MODULE ---
	profileRepo := profile.NewProfileRepository(db)
	profileSvc := profile.NewProfileService(profileRepo)
	profileCtrl := profile.NewProfileController(profileSvc)

	api := app.Group("/api")

	// Group untuk /api/auth
	auth := api.Group("/auth")

	// Public routes (tanpa token)
	auth.Post("/register-test", loginCtrl.RegisterTest)
	auth.Post("/login", loginCtrl.Login)

	// Protected route (harus pakai token)
	auth.Get("/profile", middleware.AuthMiddleware, profileCtrl.GetProfile)
}

// func SetupRoutes(app *fiber.App, db *gorm.DB) {

// 	repo := login.NewRepository(db)
// 	svc := login.NewLoginService(repo)
//  	controller := login.NewController(svc)
// 	// repo := login.NewRepository(db)
// 	// controller := login.NewController(repo)

// 	api := app.Group("/api/auth")

// 	api.Post("/register-test", controller.RegisterTest)
// 	api.Post("/login", controller.Login)

// 	profileRepo := profile.NewProfileRepository(db)
// 	profileService := profile.NewProfileService(profileRepo)
// 	profileController := profile.NewProfileController(profileService)
// 	profile := api.Group("/profil", middleware.AuthMiddleware)
// 	profile.Get("/me", profileController.GetProfile)
// 	//api.Get("/profil", middleware.AuthMiddleware, controller.Profil)
// }

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
