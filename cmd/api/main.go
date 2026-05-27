package main

import (
	"starter-wahcah-be/config"
	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Konek DB
	db := config.NewDatabase()

	// 2. Auto Migrate (Hanya di dev environment)
	db.AutoMigrate(&models.User{})

	// 3. Init App
	app := fiber.New()

	// 4. Setup Routes
	router.SetupRoutes(app, db)

	// 5. Start (Port 8080 sesuai expose dockerfile)
	app.Listen(":8080")
}
