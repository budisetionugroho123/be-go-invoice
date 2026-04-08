package models

type InvoiceDetail struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceID uint    `json:"invoice_id" gorm:"not null;index"`
	ItemID    uint    `json:"item_id" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"type:decimal(15,2);not null"`
	Subtotal  float64 `json:"subtotal" gorm:"type:decimal(15,2);not null"`

	// Relations
	Item Item `json:"item" gorm:"foreignKey:ItemID;references:ID"`
}
