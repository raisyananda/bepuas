package route

import (
	"bepuas/app/repository"
	"bepuas/app/service"
	"bepuas/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(r fiber.Router, userSvc *service.UserService, authRepo *repository.AuthRepository) {
	users := r.Group("/api/v1/users", middleware.AuthRequired(authRepo), middleware.RequireRole("Admin"))

	users.Get("/", middleware.RequirePermission("user:read"), userSvc.ListUsers)
	users.Get("/:id", middleware.RequirePermission("user:read"), userSvc.GetUserByID)
	users.Post("/", middleware.RequirePermission("user:create"), userSvc.CreateUser)
	users.Put("/:id", middleware.RequirePermission("user:update"), userSvc.UpdateUser)
	users.Delete("/:id", middleware.RequirePermission("user:delete"), userSvc.DeleteUser)
	users.Put("/:id/role", middleware.RequirePermission("user:assign-role"), userSvc.UpdateUserRole)
}
