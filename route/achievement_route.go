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
	r.Get("/", middleware.RequirePermission("achievement:read:own, achievement:read:advisee, achievement:read:all"), achievementService.List)

	// GET Detail
	r.Get("/:id", middleware.RequirePermission("achievement:read:own, achievement:read:advisee, achievement:read:all"), achievementService.Detail)

	// FR-003 Mahasiswa submit draft
	r.Post("/", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:create"), achievementService.SubmitDraft)

	// PUT Update draft
	r.Put("/:id", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:update"), achievementService.Update)

	// FR-005 Hapus draft
	r.Delete("/:id", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:delete"), refService.DeleteDraft)

	// FR-004 Submit verifikasi
	r.Post("/:id/submit", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:submit"), refService.SubmitForVerificationServ)

	r.Post("/:id/verify", middleware.RequireRole("Dosen Wali"), middleware.RequirePermission("achievement:verify"), refService.Verify)
	r.Post("/:id/reject", middleware.RequireRole("Dosen Wali"), middleware.RequirePermission("achievement:reject"), refService.Reject)
	r.Get("/:id/history", middleware.RequirePermission("achievement:read"), refService.History)

	r.Post("/:id/attachments", middleware.RequireRole("Mahasiswa"), middleware.RequirePermission("achievement:update"), achievementService.UploadAttachment)
}
