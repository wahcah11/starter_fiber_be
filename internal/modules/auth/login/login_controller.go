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
	var req RegisterRequest // Gunakan RegisterRequest, bukan LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// Validasi request untuk registrasi
	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	// Mendaftarkan pengguna baru
	err := c.service.RegisterUser(req) // Gunakan RegisterRequest pada service
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Hanya mengirimkan pesan sukses tanpa nama lengkap
	return ctx.JSON(fiber.Map{
		"message": "User created successfully", // Response hanya berisi pesan sukses
	})
}
