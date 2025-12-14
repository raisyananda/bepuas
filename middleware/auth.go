package middleware

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"bepuas/app/model"
	"bepuas/database"
	"bepuas/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "token tidak ditemukan"})
		}

		// Format "Bearer <token>"
		if strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := utils.ValidateToken(authHeader)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "token tidak valid"})
		}

		user := model.User{
			ID:       claims.UserID,
			Username: claims.Username,
			RoleID:   claims.RoleID,
		}
		c.Locals("user", user)

		// Ambil user dari Postgres
		var u model.User

		var is_active bool
		query := `SELECT id, username, email, password_hash, full_name, role_id, is_active
				  FROM users WHERE id = $1 AND is_active = true LIMIT 1`
		err = database.ConnectPostgres().QueryRow(query, claims.UserID).Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &is_active,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "user tidak ditemukan atau non-aktif"})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "gagal membaca user"})
		}

		c.Locals("user", u)
		c.Locals("role_id", claims.RoleID)
		c.Locals("user_id", claims.UserID)
		c.Locals("token_issued_at", time.Now())

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "akses hanya untuk admin"})
		}
		return c.Next()
	}
}
