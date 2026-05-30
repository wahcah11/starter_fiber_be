package crud

import (
	"starter-wahcah-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewCrudRepository(db)
	svc := NewCrudService(repo)
	ctrl := NewCrudController(svc)

	roles := router.Group("/roles")
	roles.Get("/", middleware.HasPermission(db, "role:read"), ctrl.GetAll)
	roles.Get("/:id", middleware.HasPermission(db, "role:read"), ctrl.GetByID)
	roles.Post("/", middleware.HasPermission(db, "role:create"), ctrl.Create)
	roles.Put("/:id", middleware.HasPermission(db, "role:update"), ctrl.Update)
	roles.Delete("/:id", middleware.HasPermission(db, "role:delete"), ctrl.Delete)
}
