package handler

import (
	"github.com/budisetionugroho123/be-go-invoice/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ItemHandler struct {
	itemService *service.ItemService
}

func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

// SearchByCode handles GET /api/items?code={kode}
// This endpoint is PUBLIC (no JWT required) to serve debounce requests.
func (h *ItemHandler) SearchByCode(c *fiber.Ctx) error {
	code := c.Query("code", "")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'code' is required",
		})
	}

	items, err := h.itemService.SearchByCode(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search items",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Items retrieved successfully",
		"data":    items,
	})
}
