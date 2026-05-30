package crud_user_role

import (
	"starter-wahcah-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := NewCrudUserRoleRepository(db)
	svc := NewCrudUserRoleService(repo)
	ctrl := NewCrudUserRoleController(svc)

	userRoles := router.Group("/user-roles")
	userRoles.Get("/", middleware.HasPermission(db, "user_role:read"), ctrl.GetAll)
	userRoles.Get("/:id", middleware.HasPermission(db, "user_role:read"), ctrl.GetByID)
	userRoles.Post("/", middleware.HasPermission(db, "user_role:create"), ctrl.Create)
	userRoles.Delete("/:id", middleware.HasPermission(db, "user_role:delete"), ctrl.Delete)
}
