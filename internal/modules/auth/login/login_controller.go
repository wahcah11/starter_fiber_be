package login

import (
	//"starter-wahcah-be/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	//"strconv"
	"time"
)

type Controller struct {
	service Service
}

// func NewLoginController(service Service) *Controller {
// 	return &Controller{service}
// }

func NewController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) RegisterTest(ctx *fiber.Ctx) error {
    var req RegisterRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }

    // Simpan langsung semua field (first_name, last_name, email, password)
    if err := c.service.Register(req); err != nil {
        return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return ctx.JSON(fiber.Map{
        "message": "Test user created",
        "first_name": req.FirstName,
        "last_name": req.LastName,
        "email": req.Email,
    })
}


// func (c *Controller) Register(ctx *fiber.Ctx) error {
// 	var req RegisterRequest
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	if err := c.service.Register(req); err != nil {
// 		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return ctx.JSON(RegisterResponse{Message: "register success"})
// }

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	user, err := c.service.Login(req)
	if err != nil{
		return ctx.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	t, _ := token.SignedString([]byte(secret))

	return ctx.JSON(LoginResponse{
		Token: t,
		FirstName: user.FirstName,
		LastName: user.LastName,
	})
	
}

// func (c *Controller) Login(ctx *fiber.Ctx) error {
// 	var req LoginRequest
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
// 	}

// 	user, err := c.service.Login(req)
// 	if err != nil {
// 		return ctx.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
// 	}

// 	claims := jwt.MapClaims{
// 		"user_id": user.ID,
// 		"exp":     time.Now().Add(time.Hour * 24).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	secret := os.Getenv("JWT_SECRET")
// 	t, _ := token.SignedString([]byte(secret))

// 	return ctx.JSON(LoginResponse{
// 		Token:     t,
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 	})
// }
// 	if errs := util.ValidateStruct(req); errs != nil {
// 		return ctx.Status(400).JSON(fiber.Map{"validation": errs})
// 	}

// 	res, err := c.service.Authenticate(req)
// 	if err != nil {
// 		return ctx.Status(401).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return ctx.JSON(fiber.Map{"data": res})
// }

// Endpoint tambahan buat bikin user pertama kali (biar bisa tes login)
// func (c *Controller) RegisterTest(ctx *fiber.Ctx) error {
// 	var req LoginRequest
// 	ctx.BodyParser(&req)
// 	c.service.RegisterUser(req.Email, req.Password)
// 	return ctx.JSON(fiber.Map{"message": "User created"})
// }
