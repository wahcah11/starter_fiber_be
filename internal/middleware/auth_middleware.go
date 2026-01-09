package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	//"starter_fiber_be/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

// func AuthMiddleware(ctx *fiber.Ctx) error {
//     authHeader := ctx.Get("Authorization")
//     if authHeader == "" {
//         return ctx.Status(401).JSON(fiber.Map{"error": "missing token"})
//     }

//     parts := strings.Split(authHeader, " ")
//     if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid header format"})
//     }

//     tokenString := parts[1]

//     token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
//         return []byte(os.Getenv("JWT_SECRET")), nil
//     })
//     if err != nil || !token.Valid {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid token"})
//     }

//     claims, ok := token.Claims.(jwt.MapClaims)
//     if !ok {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid claims"})
//     }

//     userIDFloat, ok := claims["user_id"].(float64)
//     if !ok {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid user_id"})
//     }

//     ctx.Locals("user_id", uint(userIDFloat))
//     return ctx.Next()
// }


// func AuthMiddleware(ctx *fiber.Ctx) error {
//     authHeader := ctx.Get("Authorization")
//     if authHeader == "" {
//         return ctx.Status(401).JSON(fiber.Map{"error": "missing token"})
//     }

//     // Format harus "Bearer <token>"
//     parts := strings.Split(authHeader, " ")
//     if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid token format"})
//     }

//     tokenStr := parts[1]

//     token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
//         return []byte(os.Getenv("JWT_SECRET")), nil
//     })

//     if err != nil || !token.Valid {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid token"})
//     }

//     claims, ok := token.Claims.(jwt.MapClaims)
//     if !ok {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid claims"})
//     }

//     // Pastikan user_id float64 (JWT decode)
//     userID, ok := claims["user_id"].(float64)
//     if !ok {
//         return ctx.Status(401).JSON(fiber.Map{"error": "invalid user_id"})
//     }

//     ctx.Locals("user_id", uint(userID))
//     return ctx.Next()
// }



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
	userIDFloat := claims["user_id"].(float64)

	ctx.Locals("user_id", uint(userIDFloat))
	// claims := token.Claims.(jwt.MapClaims)
	// userID := claims["user_id"]

	// ctx.Locals("user_id", userID)
	return ctx.Next()
}

// func Protected() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
// 		}

// 		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("JWT_SECRET")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			return c.Status(401).JSON(fiber.Map{"error": "Invalid Token"})
// 		}

// 		claims := token.Claims.(jwt.MapClaims)
// 		// Simpan user_id ke Locals agar bisa dipakai di Controller
// 		c.Locals("user_id", uint(claims["user_id"].(float64)))

// 		return c.Next()
// 	}
// }
