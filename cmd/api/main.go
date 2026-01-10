package main

import (
	"starter-wahcah-be/config"
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/router"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	
)

func main() {
	godotenv.Load()
	// 1. Konek DB
	db := config.NewDatabase()
//	if db == nil {
//	log.Println("Running without database connection")
//}

	// 2. Auto Migrate (Hanya di dev environment)
	db.AutoMigrate(&login.User{})

	// 3. Init App
	app := fiber.New()

	// 4. Setup Routes
	router.SetupRoutes(app, db)

	// 5. Start (Port 8080 sesuai expose dockerfile)
	app.Listen(":8080")

	
}
