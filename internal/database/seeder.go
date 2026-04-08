package database

import (
	"fmt"
	"log"

	"github.com/budisetionugroho123/be-go-invoice/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedUsers(db)
	seedItems(db)
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("⏩ Users already seeded, skipping...")
		return
	}

	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	keraniPass, _ := bcrypt.GenerateFromPassword([]byte("kerani123"), bcrypt.DefaultCost)

	users := []models.User{
		{Username: "admin", Password: string(adminPass), Role: "Admin"},
		{Username: "kerani", Password: string(keraniPass), Role: "Kerani"},
	}

	if err := db.Create(&users).Error; err != nil {
		log.Fatalf("❌ Failed to seed users: %v", err)
	}

	fmt.Println("✅ Users seeded successfully")
}

func seedItems(db *gorm.DB) {
	var count int64
	db.Model(&models.Item{}).Count(&count)
	if count > 0 {
		fmt.Println("⏩ Items already seeded, skipping...")
		return
	}

	items := []models.Item{
		{Code: "BRG-001", Name: "Beras Premium 5kg", Price: 75000},
		{Code: "BRG-002", Name: "Gula Pasir 1kg", Price: 14500},
		{Code: "BRG-003", Name: "Minyak Goreng 2L", Price: 32000},
		{Code: "BRG-004", Name: "Tepung Terigu 1kg", Price: 12000},
		{Code: "BRG-005", Name: "Kopi Bubuk 250g", Price: 25000},
		{Code: "BRG-006", Name: "Teh Celup 25pcs", Price: 9500},
		{Code: "BRG-007", Name: "Susu UHT 1L", Price: 18000},
		{Code: "BRG-008", Name: "Sabun Mandi 100g", Price: 5500},
		{Code: "BRG-009", Name: "Sampo 170ml", Price: 22000},
		{Code: "BRG-010", Name: "Detergen 1kg", Price: 28000},
	}

	if err := db.Create(&items).Error; err != nil {
		log.Fatalf("❌ Failed to seed items: %v", err)
	}

	fmt.Println("✅ Items seeded successfully")
}
