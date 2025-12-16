package service

import (
	"bepuas/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	studentRepo *repository.StudentRepository
	refRepo     *repository.AchievementRefRepository
}

func NewStudentService(
	studentRepo *repository.StudentRepository,
	refRepo *repository.AchievementRefRepository,
) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		refRepo:     refRepo,
	}
}

func (s *StudentService) ListStudents(c *fiber.Ctx) error {
	data, err := s.studentRepo.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (s *StudentService) GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	st, err := s.studentRepo.FindByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "student tidak ditemukan"})
	}
	return c.JSON(st)
}

func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")
	role := c.Locals("role").(string)

	// Mahasiswa hanya boleh lihat prestasi sendiri
	if role == "Mahasiswa" {
		if sid, ok := c.Locals("student_id").(string); ok && sid != studentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	data, err := s.refRepo.FindByStudentID(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}

func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		AdvisorID string `json:"advisor_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	if err := s.studentRepo.UpdateAdvisor(id, req.AdvisorID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "advisor assigned",
	})
}
