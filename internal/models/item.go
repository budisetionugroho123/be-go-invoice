package models

import "time"

type Item struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(20);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Price     float64   `json:"price" gorm:"type:decimal(15,2);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
