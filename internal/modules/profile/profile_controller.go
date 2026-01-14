package profile

import (
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewProfileController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id")

	var id uint

	if v, ok := userID.(float64); ok {
		id = uint(v)
	}

	if v, ok := userID.(uint); ok {
		id = v
	}

	if id == 0 {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	user, err := c.service.GetByID(id)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(fiber.Map{
		"full_name": user.FirstName + " " + user.LastName,
		"email":     user.Email,
	})
}

// func (c *Controller) Profile(ctx *fiber.Ctx) error {
// 	userID := ctx.Locals("user_id")
// 	if userID == nil {
// 		return ctx.Status(400).JSON(fiber.Map{"error": "user_id missing"})
// 	}

// 	id := uint(userID.(float64))

// 	user, err := c.service.GetByID(id)
// 	if err != nil {
// 		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
// 	}

// 	return ctx.JSON(fiber.Map{
// 		"first_name": user.FirstName,
// 		"last_name":  user.LastName,
// 		"email":      user.Email,
// 	})
// }
