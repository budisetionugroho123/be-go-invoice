package handler

import (
	"github.com/budisetionugroho123/be-go-invoice/internal/service"
	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	invoiceService *service.InvoiceService
}

func NewInvoiceHandler(invoiceService *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{invoiceService: invoiceService}
}

// CreateInvoice handles POST /api/invoices (JWT-protected)
func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var req service.CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get authenticated user ID from JWT middleware
	userID := c.Locals("user_id").(uint)

	invoice, err := h.invoiceService.CreateInvoice(req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Invoice created successfully",
		"data":    invoice,
	})
}
