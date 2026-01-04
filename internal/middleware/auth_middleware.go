package middleware

import (
	"fmt"
	"log"
	"os"
	"starter-wahcah-be/internal/util"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid Token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		// Simpan user_id ke Locals agar bisa dipakai di Controller
		c.Locals("user_id", uint(claims["user_id"].(float64)))
		log.Print("User id : ", uint(claims["user_id"].(float64)))

		return c.Next()
	}
}

// JWR Middleware

func JWTMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(401).JSON(fiber.Map{"error": "Missing token"})
	}

	tokenStr :=strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := util.ParseToken(tokenStr)

	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{"error": "Token Tidak cocok"})
	}

	ctx.Locals("userID", claims.UserID)
	fmt.Println(ctx.Locals("userID"))

	return ctx.Next()
}
