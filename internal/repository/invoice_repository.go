package repository

import (
	"fmt"

	"github.com/budisetionugroho123/be-go-invoice/internal/models"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

// CreateWithTransaction performs an ACID transaction:
// 1. Insert the invoice header
// 2. Insert all invoice detail rows
// If any step fails, the entire transaction is rolled back.
func (r *InvoiceRepository) CreateWithTransaction(invoice *models.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Step 1: Insert header (without details — GORM would try to auto-create them)
		details := invoice.Details
		invoice.Details = nil

		if err := tx.Create(invoice).Error; err != nil {
			return fmt.Errorf("failed to insert invoice header: %w", err)
		}

		// Step 2: Assign InvoiceID to each detail and insert
		for i := range details {
			details[i].InvoiceID = invoice.ID
		}

		if err := tx.Create(&details).Error; err != nil {
			return fmt.Errorf("failed to insert invoice details: %w", err)
		}

		// Re-attach details to the invoice struct for the response
		invoice.Details = details

		return nil
	})
}

// FindByID returns a single invoice with its details and creator preloaded.
func (r *InvoiceRepository) FindByID(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.db.Preload("Details.Item").Preload("Creator").First(&invoice, id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}
