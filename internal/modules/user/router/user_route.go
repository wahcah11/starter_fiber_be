package route

import (
	"starter-wahcah-be/internal/middleware"
	"starter-wahcah-be/internal/modules/user/controller"
	"starter-wahcah-be/internal/modules/user/repository"
	"starter-wahcah-be/internal/modules/user/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRoutes(router fiber.Router, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	ctrl := controller.NewUserController(svc)

	userGroup := router.Group("/users", middleware.Protected()) // harus login admin
	userGroup.Post("/", ctrl.CreateUser)
	userGroup.Get("/", ctrl.GetAllUsers)
	userGroup.Get("/:id", ctrl.GetUserByID)
	userGroup.Put("/:id", ctrl.UpdateUser)
	userGroup.Delete("/:id", ctrl.DeleteUser)
}
