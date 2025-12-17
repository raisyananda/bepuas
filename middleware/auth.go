package middleware

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"bepuas/app/model"
	"bepuas/app/repository"
	"bepuas/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(authRepo *repository.AuthRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).
				JSON(fiber.Map{"error": "token tidak ditemukan"})
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := utils.ValidateToken(authHeader)
		if err != nil {
			return c.Status(http.StatusUnauthorized).
				JSON(fiber.Map{"error": "token tidak valid"})
		}

		var u model.User
		var roleName sql.NullString

		err = authRepo.DB.QueryRow(`
			SELECT u.id, u.username, u.email, u.full_name, u.is_active, r.name
			FROM users u
			LEFT JOIN roles r ON r.id = u.role_id
			WHERE u.id = $1
			LIMIT 1
		`, claims.UserID).Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.FullName,
			&u.IsActive,
			&roleName,
		)

		if err != nil {
			return c.Status(http.StatusUnauthorized).
				JSON(fiber.Map{"error": "user tidak ditemukan atau non-aktif"})
		}

		perms, err := authRepo.GetUserPermissions(u.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "gagal ambil permission",
			})
		}

		log.Println("role:", roleName)
		log.Println("permission:", perms)

		// student
		var studentID sql.NullString
		_ = authRepo.DB.QueryRow(`
			SELECT id FROM students WHERE user_id = $1
		`, u.ID).Scan(&studentID)

		// lecturer
		var lecturerID sql.NullString
		_ = authRepo.DB.QueryRow(`
			SELECT id FROM lecturers WHERE user_id = $1
		`, u.ID).Scan(&lecturerID)

		// SET CONTEXT
		c.Locals("user", u)
		c.Locals("user_id", u.ID)
		c.Locals("permissions", perms)

		if roleName.Valid {
			c.Locals("role", roleName.String)
		} else {
			c.Locals("role", "UNASSIGNED")
		}

		if studentID.Valid {
			c.Locals("student_id", studentID.String)
		}
		if lecturerID.Valid {
			c.Locals("lecturer_id", lecturerID.String)
		}

		if !u.IsActive {
			return c.Status(403).JSON(fiber.Map{
				"error": "user tidak aktif",
			})
		}

		return c.Next()
	}
}

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		r := c.Locals("role")
		if r == nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		roleName, ok := r.(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		if roleName != role {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

func RequirePermission(p string) fiber.Handler {
	required := strings.Split(p, ",")

	for i := range required {
		required[i] = strings.TrimSpace(required[i])
	}

	return func(c *fiber.Ctx) error {
		perms, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

		for _, userPerm := range perms {
			for _, req := range required {
				if userPerm == req {
					return c.Next()
				}
			}
		}

		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
}
