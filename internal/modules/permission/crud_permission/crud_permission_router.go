package crud_permission

import (
	"starter-wahcah-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewCrudPermissionRepository(db)
	svc := NewCrudPermissionService(repo)
	ctrl := NewCrudPermissionController(svc)

	permissions := router.Group("/permissions")
	permissions.Get("/", middleware.HasPermission(db, "permission:read"), ctrl.GetAll)
	permissions.Get("/:id", middleware.HasPermission(db, "permission:read"), ctrl.GetByID)
	permissions.Post("/", middleware.HasPermission(db, "permission:create"), ctrl.Create)
	permissions.Put("/:id", middleware.HasPermission(db, "permission:update"), ctrl.Update)
	permissions.Delete("/:id", middleware.HasPermission(db, "permission:delete"), ctrl.Delete)
}
