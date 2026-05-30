package router

import (
	"starter-wahcah-be/internal/modules/auth/login"
	"starter-wahcah-be/internal/modules/auth/register"
	"starter-wahcah-be/internal/modules/role/crud"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// proteksi
	// middleware.OptionalAuth() -> optional login atau tidak
	// middleware.Protected() -> wajib login
	// middleware.HasPermission(db, "role:read") -> cek permissions

	login.InitRoutes(v1, db)
	register.InitRoutes(v1, db)
	crud.InitRoutes(v1, db)
}
