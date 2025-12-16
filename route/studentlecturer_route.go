package route

import (
	"bepuas/app/repository"
	"bepuas/app/service"
	"bepuas/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentLecturerRoutes(r fiber.Router, studentService *service.StudentService, lecturerService *service.LecturerService, authRepo *repository.AuthRepository) {
	api := r.Group("/api/v1", middleware.AuthRequired(authRepo))

	// STUDENTS
	students := api.Group("/students")

	students.Get("/", middleware.RequirePermission("student:read:all"), studentService.ListStudents)
	students.Get("/:id", middleware.RequirePermission("student:read"), studentService.GetStudentByID)
	students.Get("/:id/achievements", middleware.RequirePermission("achievement:read"), studentService.GetStudentAchievements)
	students.Put("/:id/advisor", middleware.RequirePermission("student:assign-advisor"), studentService.AssignAdvisor)

	// LECTURERS
	lecturers := api.Group("/lecturers")

	lecturers.Get("/", middleware.RequirePermission("lecturer:read:all"), lecturerService.ListLecturers)
	lecturers.Get("/:id/advisees", middleware.RequirePermission("student:read:advisee"), lecturerService.GetAdvisees)
}
