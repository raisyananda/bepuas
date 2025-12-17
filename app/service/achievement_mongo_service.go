package service

import (
	"context"
	"time"

	"bepuas/app/model"
	"bepuas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementMongoService struct {
	repo        *repository.AchievementMongoRepository
	refRepo     *repository.AchievementRefRepository
	studentRepo *repository.StudentRepository
}

func NewAchievementMongoService(r *repository.AchievementMongoRepository, refRepo *repository.AchievementRefRepository, studentRepo *repository.StudentRepository,
) *AchievementMongoService {
	return &AchievementMongoService{
		repo:        r,
		refRepo:     refRepo,
		studentRepo: studentRepo,
	}
}

// Submit Prestasi (Draft)
func (s *AchievementMongoService) SubmitDraft(c *fiber.Ctx) error {
	ctx := context.Background()
	studentID := c.Locals("student_id").(string)

	var a model.Achievement
	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	a.StudentID = studentID
	a.Status = "draft"

	if err := s.repo.Create(ctx, &a); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal simpan prestasi"})
	}

	// PostgreSQL reference
	if err := s.refRepo.Create(
		uuid.NewString(),
		studentID,
		a.ID.Hex(),
		"draft",
	); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal simpan reference"})
	}

	return c.JSON(a)
}

// Dipakai oleh Ref Service
func (s *AchievementMongoService) SubmitForVerification(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.SubmitForVerification(ctx, id)
}

func (s *AchievementMongoService) DeleteDraft(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *AchievementMongoService) List(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)

	// MAHASISWA
	if role == "Mahasiswa" {
		studentID := c.Locals("student_id").(string)
		data, err := s.repo.FindByStudent(ctx, studentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	}

	// DOSEN WALI
	if role == "Dosen Wali" {
		lecturerID := c.Locals("lecturer_id").(string)

		// 1. Ambil mahasiswa bimbingan
		students, err := s.studentRepo.FindByAdvisor(lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		if len(students) == 0 {
			return c.JSON([]interface{}{})
		}

		// 2. Map student_id
		studentIDs := make(map[string]bool)
		for _, s := range students {
			studentIDs[s.ID] = true
		}

		// 3. Ambil semua reference
		refs, err := s.refRepo.FindAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		var mongoIDs []primitive.ObjectID
		refMap := map[string]model.AchievementReference{}

		for _, r := range refs {
			if !studentIDs[r.StudentID] {
				continue
			}

			oid, err := primitive.ObjectIDFromHex(r.MongoAchievementID)
			if err != nil {
				continue
			}

			mongoIDs = append(mongoIDs, oid)
			refMap[r.MongoAchievementID] = r
		}

		achievements, err := s.repo.FindByIDs(ctx, mongoIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		var result []fiber.Map
		for _, a := range achievements {
			ref := refMap[a.ID.Hex()]
			result = append(result, fiber.Map{
				"id":         a.ID.Hex(),
				"title":      a.Title,
				"type":       a.AchievementType,
				"status":     ref.Status,
				"student_id": ref.StudentID,
				"created_at": a.CreatedAt,
			})
		}

		return c.JSON(result)
	}

	// ADMIN (FR-010)
	if role == "Admin" {
		refs, err := s.refRepo.FindAll()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		var mongoIDs []primitive.ObjectID
		refMap := map[string]model.AchievementReference{}

		for _, r := range refs {
			oid, err := primitive.ObjectIDFromHex(r.MongoAchievementID)
			if err != nil {
				continue
			}
			mongoIDs = append(mongoIDs, oid)
			refMap[r.MongoAchievementID] = r
		}

		achievements, err := s.repo.FindByIDs(ctx, mongoIDs)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Gabungkan reference + mongo (response admin)
		var result []fiber.Map
		for _, a := range achievements {
			ref := refMap[a.ID.Hex()]
			result = append(result, fiber.Map{
				"id":         a.ID.Hex(),
				"title":      a.Title,
				"type":       a.AchievementType,
				"status":     ref.Status,
				"student_id": ref.StudentID,
				"created_at": a.CreatedAt,
			})
		}

		return c.JSON(result)
	}

	return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
}

func (s *AchievementMongoService) Detail(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "prestasi tidak ditemukan"})
	}

	role := c.Locals("role").(string)

	// MAHASISWA
	if role == "Mahasiswa" {
		studentID := c.Locals("student_id").(string)
		if a.StudentID != studentID {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	// DOSEN WALI
	if role == "Dosen Wali" {
		lecturerID := c.Locals("lecturer_id").(string)

		ok, err := s.studentRepo.IsAdvisee(a.StudentID, lecturerID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
	}

	return c.JSON(a)
}

func (s *AchievementMongoService) Update(c *fiber.Ctx) error {
	ctx := context.Background()
	role := c.Locals("role").(string)
	studentID := c.Locals("student_id").(string)

	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	var req struct {
		Title           string                 `json:"title"`
		AchievementType string                 `json:"achievementType"`
		Description     string                 `json:"description"`
		Details         map[string]interface{} `json:"details"`
		Tags            []string               `json:"tags"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "body tidak valid"})
	}

	update := bson.M{
		"title":           req.Title,
		"achievementType": req.AchievementType,
		"description":     req.Description,
		"details":         req.Details,
		"tags":            req.Tags,
		"updated_at":      time.Now(),
	}

	if err := s.repo.UpdateDraft(ctx, id, studentID, update); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal update prestasi"})
	}

	if role != "Mahasiswa" {
		return fiber.ErrForbidden
	}

	return c.JSON(fiber.Map{
		"achievement_id": id.Hex(),
		"status":         "updated",
	})
}

func (s *AchievementMongoService) UploadAttachment(c *fiber.Ctx) error {
	ctx := context.Background()
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "id tidak valid"})
	}

	studentID := c.Locals("student_id").(string)

	var req model.Attachment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "gagal membaca body"})
	}

	if req.FileURL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "file_url wajib"})
	}

	if req.FileName == "" {
		req.FileName = "unknown"
	}

	if req.FileType == "" {
		req.FileType = "application/octet-stream"
	}

	attachment := model.Attachment{
		FileName:   req.FileName,
		FileURL:    req.FileURL,
		FileType:   req.FileType,
		UploadedAt: time.Now(),
	}

	// Simpan attachment ke repo
	err = s.repo.AddAttachment(ctx, id, studentID, attachment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "uploaded",
		"file":   attachment,
	})
}
