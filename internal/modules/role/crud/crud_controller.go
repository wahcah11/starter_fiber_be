package crud

import (
	"strconv"

	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewCrudController(service Service) *Controller {
	return &Controller{service}
}

// GetAll menangani GET /roles?page=1&limit=10
func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	res, err := c.service.GetAll(page, limit)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// GetByID menangani GET /roles/:id
func (c *Controller) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	res, err := c.service.GetByID(uint(id))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// Create menangani POST /roles
func (c *Controller) Create(ctx *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	res, err := c.service.Create(req)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(201).JSON(res)
}

// Update menangani PUT /roles/:id
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	var req UpdateRoleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	res, err := c.service.Update(uint(id), req)
	if err != nil {
		if err.Error() == "role not found" {
			return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// Delete menangani DELETE /roles/:id
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid role ID"})
	}

	if err := c.service.Delete(uint(id)); err != nil {
		if err.Error() == "role not found" {
			return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Role deleted successfully",
	})
}
