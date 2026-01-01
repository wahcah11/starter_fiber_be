package profil

import "github.com/gofiber/fiber/v2"

type Controller struct {
	service Service
}

func NewProfilController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	res, err := c.service.GetProfile(userID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": res})
}
