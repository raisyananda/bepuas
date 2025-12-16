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
	refRepo  *repository.AchievementRefRepository
	mongoSvc *AchievementMongoService
}

func NewAchievementRefService(ref *repository.AchievementRefRepository, mongoSvc *AchievementMongoService) *AchievementRefService {
	return &AchievementRefService{
		refRepo:  ref,
		mongoSvc: mongoSvc,
	}
}

// FR-004: Submit untuk Verifikasi

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

// FR-005: Hapus Prestasi Draft
func (s *AchievementRefService) DeleteDraft(c *fiber.Ctx) error {
	ctx := context.Background()

	idParam := c.Params("id")

	mongoID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	// 1. Update PostgreSQL dulu
	if err := s.refRepo.DeleteDraft(idParam); err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Baru delete Mongo
	if err := s.mongoSvc.DeleteDraft(ctx, mongoID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"achievement_id": idParam,
		"status":         "deleted",
	})
}
