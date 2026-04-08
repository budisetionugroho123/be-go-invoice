package database

import (
	"fmt"
	"log"

	"github.com/budisetionugroho123/be-go-invoice/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Invoice{},
		&models.InvoiceDetail{},
	)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Database migrated successfully")
}
