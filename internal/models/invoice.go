package models

import "time"

type Invoice struct {
	ID              uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceNumber   string          `json:"invoice_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	SenderName      string          `json:"sender_name" gorm:"type:varchar(255);not null"`
	SenderAddress   string          `json:"sender_address" gorm:"type:text;not null"`
	ReceiverName    string          `json:"receiver_name" gorm:"type:varchar(255);not null"`
	ReceiverAddress string          `json:"receiver_address" gorm:"type:text;not null"`
	TotalAmount     float64         `json:"total_amount" gorm:"type:decimal(15,2);not null;default:0"`
	CreatedBy       uint            `json:"created_by" gorm:"not null"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Relations
	Creator User            `json:"creator" gorm:"foreignKey:CreatedBy;references:ID"`
	Details []InvoiceDetail  `json:"details" gorm:"foreignKey:InvoiceID;references:ID"`
}
