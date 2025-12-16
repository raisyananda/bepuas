package route

import (
	"bepuas/app/repository"
	"bepuas/app/service"
	"bepuas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(app *fiber.App, service *service.AuthService, authRepo *repository.AuthRepository) {
	api := app.Group("/api/v1/auth")
	api.Post("/login", service.Login)

	auth := app.Group("/api/v1/auth", middleware.AuthRequired(authRepo))

	// auth.Post("/login", service.Login)
	auth.Post("/refresh", service.RefreshToken)
	auth.Post("/logout", service.Logout)
	auth.Get("/profile", service.Profile)
}
