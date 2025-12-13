package main

import (
	"log"
	"os"

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

	// Routes
	route.AuthRoute(app, cfg.AuthService)
	route.PrestasiRoute(app, cfg.AchievementService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server berjalan di http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
