package route

import (
	"bepuas/app/repository"
	"bepuas/app/service"
	"bepuas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoute(app *fiber.App, achievementService *service.AchievementMongoService, refService *service.AchievementRefService, authRepo *repository.AuthRepository) {
	r := app.Group("/api/v1/achievements", middleware.AuthRequired(authRepo))

	// GET List (filtered by role)
	r.Get("/", middleware.RequirePermission("achievement:read"), achievementService.List)

	// GET Detail
	r.Get("/:id", middleware.RequirePermission("achievement:read"), achievementService.Detail)

	// FR-003 Mahasiswa submit draft
	r.Post("/", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:create"), achievementService.SubmitDraft)

	// PUT Update draft
	r.Put("/:id", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:update"), achievementService.Update)

	// FR-004 Submit verifikasi
	r.Post("/:id/submit", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:submit"), refService.SubmitForVerification)

	// FR-005 Hapus draft
	r.Delete("/:id", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:delete"), refService.DeleteDraft)
}
