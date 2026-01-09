package main

import (
	"starter-wahcah-be/config"
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Konek DB
	db := config.NewDatabase()

	// 2. Auto Migrate (Hanya di dev environment)
	db.AutoMigrate(&login.User{})

	// 3. Init App
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))


	// Root server
	app.Get("/", Hello)
	// 4. Setup Routes
	router.SetupRoutes(app, db)
	// 5. Start (Port 8080 sesuai expose dockerfile)
	app.Listen(":8080")
}


func Hello (ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("Back end siap cuy")
}
