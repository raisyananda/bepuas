package service

/*
import (
	"errors"

	"bepuas/app/model"
	"bepuas/app/repository"
	"bepuas/utils"
)

type UserService struct {
	repo *repository.NewUserRepository
}

func NewAuthService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Login(identifier, password string) (string, error) {
	user, err := s.repo.FindByUsernameOrEmail(identifier)
	if err != nil {
		return "", err
	}

	if !utils.CheckPassword(password, user.Password) {
		return "", errors.New("password salah")
	}

	token, err := utils.GenerateJWT(user.ID, user.RoleID)
	return token, err
}
*/

/*
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
*/

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

	user, err := s.repo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "user tidak ditemukan"})
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "password salah"})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{"error": "user tidak aktif"})
	}

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
	user := c.Locals("user").(model.User)

	perms, err := s.repo.GetUserPermissions(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal ambil permission"})
	}

	token, err := utils.GenerateToken(user, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "logout berhasil"})
}

func (s *AuthService) Profile(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	return c.JSON(user)
}
