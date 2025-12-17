package service

import (
	"context"
	"log"

	"bepuas/app/model"
	"bepuas/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementRefService struct {
	refRepo     *repository.AchievementRefRepository
	mongoSvc    *AchievementMongoService
	studentRepo *repository.StudentRepository
}

func NewAchievementRefService(ref *repository.AchievementRefRepository, mongoSvc *AchievementMongoService,
	studentRepo *repository.StudentRepository) *AchievementRefService {
	return &AchievementRefService{
		refRepo:     ref,
		mongoSvc:    mongoSvc,
		studentRepo: studentRepo,
	}
}

// Submit untuk Verifikasi
func (s *AchievementRefService) SubmitForVerification(c *fiber.Ctx) error {
	ctx := context.Background()
	studentID := c.Locals("student_id").(model.Student)

	idParam := c.Params("id")
	mongoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	// Mongo
	if err := s.mongoSvc.SubmitForVerification(ctx, mongoID); err != nil {
		log.Println(" MONGO SUBMIT ERROR:", err)
		return c.Status(500).JSON(fiber.Map{"error": "gagal update prestasi"})
	}

	// PostgreSQL
	if err := s.refRepo.UpdateStatus(idParam, "submitted"); err != nil {
		log.Println("PG UPDATE STATUS ERROR:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Println(" SUBMITTED:", idParam)

	return c.JSON(fiber.Map{
		"achievement_id": idParam,
		"status":         "submitted",
		"submitted_by":   studentID.ID,
	})
}

func (s *AchievementRefService) SubmitForVerificationServ(c *fiber.Ctx) error {
	ctx := context.Background()

	idParam := c.Params("id")
	mongoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	studentID := c.Locals("student_id").(string)

	// Cek ownership
	a, err := s.mongoSvc.repo.FindByID(ctx, mongoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "prestasi tidak ditemukan"})
	}
	if a.StudentID != studentID {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	//Update Postgre
	if err := s.refRepo.UpdateStatus(idParam, "submitted"); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal update reference"})
	}

	// update Mongo
	if err := s.mongoSvc.repo.SubmitForVerification(ctx, mongoID); err != nil {
		// rollback Postgre kalau mau ideal
		return c.Status(500).JSON(fiber.Map{"error": "gagal submit mongo"})
	}

	return c.JSON(fiber.Map{
		"achievement_id": idParam,
		"status":         "submitted",
	})
}

// FR-005: Hapus Prestasi Draft
func (s *AchievementRefService) DeleteDraft(c *fiber.Ctx) error {
	ctx := context.Background()

	idParam := c.Params("id")

	mongoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	// Update PostgreSQL
	if err := s.refRepo.DeleteDraft(idParam); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	// delete Mongo
	if err := s.mongoSvc.DeleteDraft(ctx, mongoID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"achievement_id": idParam,
		"status":         "deleted",
	})
}

func (s *AchievementRefService) Verify(c *fiber.Ctx) error {
	idParam := c.Params("id")
	lecturerID := c.Locals("lecturer_id").(string)
	userID := c.Locals("user_id").(string)

	ref, err := s.refRepo.FindByMongoID(idParam)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference tidak ditemukan"})
	}

	// cek advisee
	ok, err := s.studentRepo.IsAdvisee(ref.StudentID, lecturerID)
	if err != nil || !ok {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	if err := s.refRepo.Verify(idParam, userID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"achievement_id": idParam,
		"status":         "verified",
	})
}

func (s *AchievementRefService) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	lecturerID := c.Locals("lecturer_id").(string)

	var req struct {
		RejectionNote string `json:"rejection_note"`
	}
	_ = c.BodyParser(&req)

	ref, err := s.refRepo.FindByMongoID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "reference tidak ditemukan"})
	}

	ok, _ := s.studentRepo.IsAdvisee(ref.StudentID, lecturerID)
	if !ok {
		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}

	err = s.refRepo.Reject(id, req.RejectionNote)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"achievement_id": id,
		"status":         "rejected",
	})
}

func (s *AchievementRefService) History(c *fiber.Ctx) error {
	id := c.Params("id")

	history, err := s.refRepo.FindHistoryByMongoID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "history tidak ditemukan"})
	}

	return c.JSON(history)
}
