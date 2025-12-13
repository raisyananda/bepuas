package route

import (
	"bepuas/app/service"
	"bepuas/middleware"

	"github.com/gofiber/fiber/v2"
)

/*
func UserRoute(app *fiber.App, service *service.AuthService) {
	user := app.Group("/api/users")

	// Public
	user.Post("/login", service.Login)

	// Protected
	user.Get("/profile", middleware.AuthRequired(), service.GetProfile)

	// Admin Only
	// user.Get("/", middleware.AuthRequired(), middleware.AdminOnly(), service.GetAllUser)
	// user.Delete("/:id", middleware.AuthRequired(), middleware.AdminOnly(), service.DeleteUser)
}
*/

func AuthRoute(app *fiber.App, service *service.AuthService) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", service.Login)
	auth.Post("/refresh", service.RefreshToken)
	auth.Post("/logout", middleware.AuthRequired(), service.Logout)
	auth.Get("/profile", middleware.AuthRequired(), service.Profile)
}
