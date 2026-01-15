package login

import (
	"github.com/gofiber/fiber/v2"
	//"github.com/golang-jwt/jwt/v5"
	//"os"
	//"time"
)

type Controller struct {
	service Service
}

func NewLoginController(service Service) *Controller {
	return &Controller{service}
}

func (c *Controller) RegisterTest(ctx *fiber.Ctx) error {
    var req RegisterRequest
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }

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


func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	res, err := c.service.Authenticate(req)
    if err != nil {
        return ctx.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
    }

    return ctx.JSON(res)



	// user, err := c.service.Authenticate(req)
	// if err != nil{
	// 	return ctx.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	// }

	// claims := jwt.MapClaims{
	// 	"user_id": user.ID,
	// 	"exp": time.Now().Add(time.Hour * 24).Unix(),
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// secret := os.Getenv("JWT_SECRET")
	// t, _ := token.SignedString([]byte(secret))

	// return ctx.JSON(LoginResponse{
	// 	Token: t,
	// 	FirstName: user.FirstName,
	// 	LastName: user.LastName,
	// })
	
}

// func (c *Controller) Profile(ctx *fiber.Ctx) error {
// 	userID := ctx.Locals("user_id")
// 	if userID == nil {
// 		return ctx.Status(400).JSON(fiber.Map{"error": "user_id missing"})
// 	}

// 	id := uint(userID.(float64))

// 	user, err := c.service.GetByID(id)
// 	if err != nil {
// 		return ctx.Status(404).JSON(fiber.Map{"error": "User not found"})
// 	}

// 	return ctx.JSON(fiber.Map{
// 		"first_name": user.FirstName,
// 		"last_name":  user.LastName,
// 		"email":      user.Email,
// 	})
// }


