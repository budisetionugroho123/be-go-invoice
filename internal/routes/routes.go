package routes

import (
	"github.com/budisetionugroho123/be-go-invoice/internal/config"
	"github.com/budisetionugroho123/be-go-invoice/internal/handler"
	"github.com/budisetionugroho123/be-go-invoice/internal/middleware"
	"github.com/budisetionugroho123/be-go-invoice/internal/repository"
	"github.com/budisetionugroho123/be-go-invoice/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// ── Repositories ──────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)

	// ── Services ──────────────────────────────────────────────
	authService := service.NewAuthService(userRepo, cfg)
	itemService := service.NewItemService(itemRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo, itemRepo)

	// ── Handlers ──────────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authService)
	itemHandler := handler.NewItemHandler(itemService)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	// ── API Routes ────────────────────────────────────────────
	api := app.Group("/api")

	// Public routes (no JWT required)
	api.Post("/login", authHandler.Login)
	api.Get("/items", itemHandler.SearchByCode)

	// Protected routes (JWT required)
	protected := api.Group("", middleware.JWTMiddleware(cfg))
	protected.Post("/invoices", invoiceHandler.CreateInvoice)
}
