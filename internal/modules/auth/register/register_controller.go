package register

import (
	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewRegisterController(service Service) *Controller {
	return &Controller{service}
}

// Register menangani POST /auth/register
func (c *Controller) Register(ctx *fiber.Ctx) error {
	// 1. Parse request body
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// 2. Validasi struct
	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	// 3. Panggil service
	res, err := c.service.Register(req)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// 4. Return response
	return ctx.Status(201).JSON(fiber.Map{"data": res})
}
