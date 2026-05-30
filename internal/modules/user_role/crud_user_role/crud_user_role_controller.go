package crud_user_role

import (
	"strconv"

	"starter-wahcah-be/internal/util"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewCrudUserRoleController(service Service) *Controller {
	return &Controller{service}
}

// GetAll menangani GET /user-roles?page=1&limit=10&user_id=1
func (c *Controller) GetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Filter opsional by user_id
	var userID *uint
	if userIDStr := ctx.Query("user_id", ""); userIDStr != "" {
		parsed, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": "Invalid user_id"})
		}
		userIDVal := uint(parsed)
		userID = &userIDVal
	}

	res, err := c.service.GetAll(page, limit, userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// GetByID menangani GET /user-roles/:id
func (c *Controller) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid user role ID"})
	}

	res, err := c.service.GetByID(uint(id))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(res)
}

// Create menangani POST /user-roles
func (c *Controller) Create(ctx *fiber.Ctx) error {
	var req CreateUserRoleRequest
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

// Delete menangani DELETE /user-roles/:id
func (c *Controller) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil || id < 1 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid user role ID"})
	}

	if err := c.service.Delete(uint(id)); err != nil {
		if err.Error() == "user role not found" {
			return ctx.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Role removed from user successfully",
	})
}
