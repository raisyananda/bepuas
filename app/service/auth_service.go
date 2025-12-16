package service

import (
	"bepuas/app/model"
	"bepuas/app/repository"
	"bepuas/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(r *repository.AuthRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	if req.Identifier == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "identifier dan password wajib"})
	}

	// Cari user by username ATAU email
	user, err := s.repo.FindByUsernameOrEmail(req.Identifier)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "user tidak ditemukan"})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{"error": "user tidak aktif"})
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "password salah"})
	}

	// Ambil permission
	perms, err := s.repo.GetUserPermissions(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal ambil permission"})
	}

	token, err := utils.GenerateToken(user, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal generate token"})
	}

	return c.JSON(model.LoginResponse{
		User:        user,
		Token:       token,
		Permissions: perms,
	})
}

func (s *AuthService) RefreshToken(c *fiber.Ctx) error {
	userAny := c.Locals("user")
	if userAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user := userAny.(model.User)

	perms, err := s.repo.GetUserPermissions(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal ambil permission"})
	}

	token, err := utils.GenerateToken(user, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logout berhasil"})
}

func (s *AuthService) Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	return c.JSON(user)
}
