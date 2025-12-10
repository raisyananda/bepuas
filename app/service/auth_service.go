package service

import (
	"bepuas/app/model"
	"bepuas/app/repository"
	"bepuas/utils"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	userRepo *repository.AuthRepository
}

func NewAuthService(u *repository.AuthRepository) *AuthService {
	return &AuthService{userRepo: u}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	user, hash, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "not found" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	if !utils.CheckPassword(req.Password, hash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if !user.IsActive {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user inactive"})
	}
	perms, err := s.userRepo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed load permissions"})
	}
	token, err := utils.GenerateToken(user, perms)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "generate token failed"})
	}
	return c.JSON(model.LoginResponse{User: user, Token: token, Permissions: perms})
}

func (s *AuthService) GetProfile(c *fiber.Ctx) error {
	claims := c.Locals("jwt_claims").(*model.JWTClaims)
	user := model.User{
		ID:       claims.UserID,
		Username: claims.Username,
		RoleID:   claims.RoleID,
	}
	return c.JSON(fiber.Map{"user": user, "permissions": claims.Permissions})
}
