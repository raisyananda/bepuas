package main

/*
import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"p4_ProjectGo/app/repository"
	"p4_ProjectGo/app/service"
	"p4_ProjectGo/database"
	"p4_ProjectGo/route"
)

func main() {
	// Connect MongoDB
	client, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer client.Disconnect(context.Background())

	// Setup Repository
	collection := client.Database("ecommerce").Collection("products")
	repo := repository.NewProductRepository(collection)

	// Setup Service
	productService := service.NewProductService(repo)

	// Setup Fiber
	app := fiber.New()

	// Register Routes
	route.ProductRoute(app, productService)

	// Start Server
	log.Println("Server berjalan di http://localhost:3000")
	app.Listen(":3000")
}
*/
