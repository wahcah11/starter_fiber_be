package controller

import (
	"starter-wahcah-be/internal/modules/user/dto"
	"starter-wahcah-be/internal/modules/user/model"
	"starter-wahcah-be/internal/modules/user/service"
	"starter-wahcah-be/internal/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service service.UserService
}

func NewUserController(svc service.UserService) *UserController {
	return &UserController{service: svc}
}


func (c *UserController) CreateUser(ctx *fiber.Ctx) error {

	var req dto.CreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// 2. Validasi input
	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	// 4. Buat user
	user := model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Role: req.Role,
		Password:  req.Password,
	}

	if err := c.service.CreateUser(&user); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(201).JSON(req)
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.service.GetAllUsers()
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return ctx.JSON(user)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	var req model.User
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	req.ID = uint(id)
	if err := c.service.UpdateUser(&req); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(req)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.service.DeleteUser(uint(id)); err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(200).JSON(fiber.Map{"message": "berhasil menghapus user "})
}
