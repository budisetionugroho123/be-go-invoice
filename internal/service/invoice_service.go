package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/budisetionugroho123/be-go-invoice/internal/models"
	"github.com/budisetionugroho123/be-go-invoice/internal/repository"
)

type InvoiceService struct {
	invoiceRepo *repository.InvoiceRepository
	itemRepo    *repository.ItemRepository
}

func NewInvoiceService(invoiceRepo *repository.InvoiceRepository, itemRepo *repository.ItemRepository) *InvoiceService {
	return &InvoiceService{
		invoiceRepo: invoiceRepo,
		itemRepo:    itemRepo,
	}
}

// CreateInvoiceRequest is the payload sent from the frontend.
type CreateInvoiceRequest struct {
	SenderName      string                     `json:"sender_name"`
	SenderAddress   string                     `json:"sender_address"`
	ReceiverName    string                     `json:"receiver_name"`
	ReceiverAddress string                     `json:"receiver_address"`
	Items           []CreateInvoiceItemRequest `json:"items"`
}

type CreateInvoiceItemRequest struct {
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}

// CreateInvoice implements the Zero-Trust logic:
// 1. Validate all input
// 2. Query real prices from the Master Item table (NEVER trust frontend prices)
// 3. Recalculate subtotals and grand total server-side
// 4. Save atomically via db.Transaction()
func (s *InvoiceService) CreateInvoice(req CreateInvoiceRequest, userID uint) (*models.Invoice, error) {
	// Validate basic input
	if req.SenderName == "" || req.SenderAddress == "" {
		return nil, errors.New("sender name and address are required")
	}
	if req.ReceiverName == "" || req.ReceiverAddress == "" {
		return nil, errors.New("receiver name and address are required")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Collect all item IDs from the request
	itemIDs := make([]uint, len(req.Items))
	for i, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("item at index %d has invalid quantity", i)
		}
		itemIDs[i] = item.ItemID
	}

	// ZERO-TRUST: Query real prices from database
	masterItems, err := s.itemRepo.FindByIDs(itemIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch master items: %w", err)
	}

	// Build a lookup map: itemID -> Item
	itemMap := make(map[uint]models.Item)
	for _, item := range masterItems {
		itemMap[item.ID] = item
	}

	// Validate all requested items exist and build details with server-side pricing
	var grandTotal float64
	details := make([]models.InvoiceDetail, len(req.Items))

	for i, reqItem := range req.Items {
		masterItem, exists := itemMap[reqItem.ItemID]
		if !exists {
			return nil, fmt.Errorf("item with ID %d not found in master data", reqItem.ItemID)
		}

		// Calculate using REAL price from database (zero-trust)
		subtotal := masterItem.Price * float64(reqItem.Quantity)
		grandTotal += subtotal

		details[i] = models.InvoiceDetail{
			ItemID:   reqItem.ItemID,
			Quantity: reqItem.Quantity,
			Price:    masterItem.Price, // Real price from DB
			Subtotal: subtotal,         // Recalculated server-side
		}
	}

	// Generate invoice number
	invoiceNumber := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102-150405"), time.Now().UnixMilli()%10000)

	// Build invoice model
	invoice := &models.Invoice{
		InvoiceNumber:   invoiceNumber,
		SenderName:      req.SenderName,
		SenderAddress:   req.SenderAddress,
		ReceiverName:    req.ReceiverName,
		ReceiverAddress: req.ReceiverAddress,
		TotalAmount:     grandTotal, // Recalculated by server
		CreatedBy:       userID,
		Details:         details,
	}

	// Save with ACID transaction
	if err := s.invoiceRepo.CreateWithTransaction(invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Reload the invoice with relations for the response
	savedInvoice, err := s.invoiceRepo.FindByID(invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("invoice created but failed to reload: %w", err)
	}

	// ── BONUS: Async Webhook Integration
	go s.sendWebhook(savedInvoice)

	return savedInvoice, nil
}

func (s *InvoiceService) sendWebhook(invoice *models.Invoice) {
	webhookUrl := "https://webhook.site/dummy-fleetify-webhook" // Ganti dummy url ini jika untuk testing live

	payload, err := json.Marshal(invoice)
	if err != nil {
		log.Printf("❌ Webhook error: failed to marshal payload - %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("❌ Webhook error: failed to construct request - %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️ Webhook error: failed to send - %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ Webhook sent successfully for invoice %s\n", invoice.InvoiceNumber)
	} else {
		log.Printf("⚠️ Webhook warning: received status code %d\n", resp.StatusCode)
	}
}
