package main

import (
	"log"
	"os"

	"bepuas/app/repository"
	"bepuas/config"
	"bepuas/database"
	"bepuas/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadEnv()

	// PostgreSQL (Auth)
	pg := database.ConnectPostgres()

	// MongoDB (Achievement)
	mongo := database.ConnectMongoDB()

	cfg := config.NewConfig(pg, mongo)

	app := fiber.New()
	app.Use(logger.New())

	authRepo := repository.NewAuthRepository(pg)

	// Routes
	route.AuthRoute(app, cfg.AuthService, authRepo)
	route.AchievementRoute(app, cfg.AchievementMongoService, cfg.AchievementRefService, authRepo)
	route.UserRoutes(app, cfg.UserService, authRepo)
	route.StudentLecturerRoutes(app, cfg.StudentService, cfg.LectureService, authRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server berjalan di http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
