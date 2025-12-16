package service

import (
	"bepuas/app/model"
	"bepuas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) ListUsers(c *fiber.Ctx) error {
	users, err := s.repo.FindAllUser()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (s *UserService) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.repo.FindUserByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user tidak ditemukan"})
	}
	return c.JSON(user)
}

func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal hash password"})
	}

	user := model.User{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		IsActive: true,
	}

	if err := s.repo.CreateUser(&user, string(hash)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req model.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	// Ambil role user
	roleID, roleName, err := s.repo.GetUserRole(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "gagal mengambil role user",
		})
	}

	if roleID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "role user belum diset",
		})
	}

	// ===== LECTURER =====
	if roleName == "Dosen Wali" {
		if req.Lecturer == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "data lecturer wajib diisi",
			})
		}

		err := s.repo.UpsertLecturer(&model.Lecturer{
			ID:         uuid.NewString(),
			UserID:     userID,
			LecturerID: req.Lecturer.LecturerID,
			Department: req.Lecturer.Department,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"status": "lecturer profile updated",
		})
	}

	// ===== STUDENT =====
	if roleName == "Mahasiswa" {
		if req.Student == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "data student wajib diisi",
			})
		}

		err := s.repo.UpsertStudent(&model.Student{
			ID:           uuid.NewString(),
			UserID:       userID,
			StudentID:    req.Student.StudentID,
			ProgramStudy: req.Student.ProgramStudy,
			AcademicYear: req.Student.AcademicYear,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"status": "student profile updated",
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"error": "role tidak didukung",
	})
}

func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "deleted"})
}

func (s *UserService) UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	if err := s.repo.UpdateUserRole(id, req.RoleID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "role updated"})
}
