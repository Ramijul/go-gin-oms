package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	OrderID          uuid.UUID `gorm:"not null"`
	Order            Order     `gorm:"primaryKey;foreignKey:OrderID;references:ID"`
	ProductID        uuid.UUID `gorm:"not null"`
	Product          Product   `gorm:"primaryKey;foreignKey:ProductID;references:ID"`
	ProductName      string    `gorm:"not null"`
	ProductUnitPrice float64   `gorm:"type:numeric(8,2);not null"`
	Units            int       `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime;not null"`
}
