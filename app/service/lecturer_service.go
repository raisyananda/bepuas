package service

import (
	"bepuas/app/repository"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	lecturerRepo *repository.LecturerRepository
	studentRepo  *repository.StudentRepository
}

func NewLecturerService(lecturerRepo *repository.LecturerRepository, studentRepo *repository.StudentRepository,
) *LecturerService {
	return &LecturerService{
		lecturerRepo: lecturerRepo,
		studentRepo:  studentRepo,
	}
}

func (s *LecturerService) ListLecturers(c *fiber.Ctx) error {
	data, err := s.lecturerRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")
	role := c.Locals("role").(string)

	// Dosen hanya boleh lihat advisee sendiri
	if role == "Dosen Wali" {
		if lid, ok := c.Locals("lecturer_id").(string); ok && lid != lecturerID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	data, err := s.studentRepo.FindByAdvisor(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}
