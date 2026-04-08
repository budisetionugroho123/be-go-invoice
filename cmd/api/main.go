package main

import (
	"fmt"
	"log"

	"github.com/budisetionugroho123/be-go-invoice/internal/config"
	"github.com/budisetionugroho123/be-go-invoice/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect to database
	db := database.Connect(cfg)

	// Auto-migrate & seed
	database.Migrate(db)
	database.Seed(db)

	// Initialize Fiber
	app := fiber.New(fiber.Config{
		AppName: "Invoice API v1.0",
	})

	// Global middleware
	app.Use(fiberLogger.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Invoice API is running 🚀",
		})
	})

	// TODO: Register routes here
	_ = db // will be used when registering routes

	// Start server
	port := cfg.AppPort
	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
