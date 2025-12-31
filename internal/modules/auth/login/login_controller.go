package login

import (
	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewLoginController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	res, err := c.service.Authenticate(req)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": res})
}

// Endpoint tambahan buat bikin user pertama kali (biar bisa tes login)
func (c *Controller) RegisterTest(ctx *fiber.Ctx) error {
	var req LoginRequest
	ctx.BodyParser(&req)
	c.service.RegisterUser(req.Email, req.Password)
	return ctx.JSON(fiber.Map{"message": "User created"})
}
