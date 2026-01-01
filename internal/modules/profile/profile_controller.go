package profile

import (
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewProfileController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	// Ambil userID dari JWT middleware
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok {
		return ctx.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	profile, err := c.service.GetProfile(userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.Status(200).JSON(profile)
}
