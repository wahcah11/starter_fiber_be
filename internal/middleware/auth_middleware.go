package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
    authHeader := ctx.Get("Authorization")
    if authHeader == "" {
        return ctx.Status(401).JSON(fiber.Map{"error": "missing token"})
    }

    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil || !token.Valid {
        return ctx.Status(401).JSON(fiber.Map{"error": "invalid token"})
    }

    claims := token.Claims.(jwt.MapClaims)
    ctx.Locals("user_id", uint(claims["user_id"].(float64)))

    return ctx.Next()
}
