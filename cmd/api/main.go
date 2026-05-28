package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"starter-wahcah-be/config"
	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/router"
	"starter-wahcah-be/internal/seeder"
)

func main() {
	// 1. Koneksi Database
	db := config.NewDatabase()

	// 2. Auto Migrate
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 3. Seeder (skip di production)
	if os.Getenv("APP_ENV") != "production" {
		if err := seeder.Run(db); err != nil {
			log.Fatalf("Seeder failed: %v", err)
		}
	}

	// 4. Init Fiber
	app := fiber.New(fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		ErrorHandler: errorHandler,
	})

	// 5. Global Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// 6. Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"env":    os.Getenv("APP_ENV"),
		})
	})

	// 7. Setup Routes
	router.SetupRoutes(app, db)

	// 8. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s [%s]", port, os.Getenv("APP_ENV"))
	log.Fatal(app.Listen(":" + port))
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
