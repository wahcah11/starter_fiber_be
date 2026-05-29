package middleware

import (
	"os"
	"strings"

	"starter-wahcah-be/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// HasPermission menerima satu atau lebih permission.
// Jika lebih dari satu, user harus memiliki SEMUA permission yang disebutkan.
//
// Contoh penggunaan:
//
//	middleware.HasPermission(db, "role:read")
//	middleware.HasPermission(db, "role:create", "role:read")
func HasPermission(db *gorm.DB, permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil token dari header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		// 2. Ambil user_id dari claims
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		// 3. Load semua permissions milik user dari DB
		var userRoles []models.UserRole
		if err := db.
			Preload("Role.Permissions").
			Where("user_id = ?", userID).
			Find(&userRoles).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to load permissions"})
		}

		// 4. Kumpulkan semua permission name milik user ke dalam map
		userPermissions := make(map[string]bool)
		for _, ur := range userRoles {
			for _, perm := range ur.Role.Permissions {
				userPermissions[perm.Name] = true
			}
		}

		// 5. Cek apakah user memiliki SEMUA permission yang dibutuhkan
		for _, required := range permissions {
			if !userPermissions[required] {
				return c.Status(403).JSON(fiber.Map{
					"error":    "Forbidden",
					"required": permissions,
				})
			}
		}

		// 6. Simpan user_id ke Locals untuk dipakai controller
		c.Locals("user_id", userID)

		return c.Next()
	}
}
