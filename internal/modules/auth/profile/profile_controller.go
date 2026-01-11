package profile

import "github.com/gofiber/fiber/v2"

type Controller struct {
	service Service
}

func NewProfileController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	userIDAny := ctx.Locals("user_id") // ✅ KEY BENAR
	if userIDAny == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	userID := userIDAny.(uint)

	profile, err := c.service.GetProfile(userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(profile)
}
