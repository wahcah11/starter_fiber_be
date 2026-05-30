package crud_permission

import (
	"strconv"

	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewCrudPermissionController(service Service) *Controller {
	return &Controller{service}
}

// GetAll menangani GET /permissions?page=1&limit=10&role_id=1
func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Filter opsional by role_id
	var roleID *uint
	if roleIDStr := ctx.Query("role_id", ""); roleIDStr != "" {
		parsed, err := strconv.ParseUint(roleIDStr, 10, 32)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": "Invalid role_id"})
		}
		roleIDVal := uint(parsed)
		roleID = &roleIDVal
	}

	res, err := c.service.GetAll(page, limit, roleID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// GetByID menangani GET /permissions/:id
func (c *Controller) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid permission ID"})
	}

	res, err := c.service.GetByID(uint(id))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// Create menangani POST /permissions
func (c *Controller) Create(ctx *fiber.Ctx) error {
	var req CreatePermissionRequest
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

// Update menangani PUT /permissions/:id
func (c *Controller) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid permission ID"})
	}

	var req UpdatePermissionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	res, err := c.service.Update(uint(id), req)
	if err != nil {
		if err.Error() == "permission not found" {
			return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// Delete menangani DELETE /permissions/:id
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid permission ID"})
	}

	if err := c.service.Delete(uint(id)); err != nil {
		if err.Error() == "permission not found" {
			return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Permission deleted successfully",
	})
}
