package login

import (
	"starter-wahcah-be/internal/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
	db *gorm.DB
}

func NewLoginController(service Service, db *gorm.DB) *Controller {
	return &Controller{service: service, db: db,}
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

	return ctx.JSON(res)
}

// Endpoint tambahan buat bikin user pertama kali (biar bisa tes login)
func (c *Controller) Register(ctx *fiber.Ctx) error {
	// 1. Parse request
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	// 2. Validasi input
	if errs := util.ValidateStruct(req); errs != nil {
		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// 4. Buat user
	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	if err := c.db.Create(&user).Error; err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
	}


	// 7. Return response
	return ctx.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"id":        user.ID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"createdAt": user.CreatedAt.Format(time.RFC3339),
		},
	})
}

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uint)

	if(!ok) {
		return ctx.Status(401).JSON(fiber.Map{"error": "Unautorized"})

	}
	var user User

	if err := c.db.First(&user, userID).Error; err != nil {
		return ctx.Status(404).JSON(fiber.Map{"error": "User not Found"})


	}

	fullName := user.FirstName + " " + user.LastName

	res := ProfileResponse {
		Name: fullName,
		Email: user.Email,
	}

	return ctx.Status(200).JSON(res)
}

