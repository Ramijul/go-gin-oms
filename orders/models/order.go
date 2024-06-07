package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	Street   string `gorm:"not null"`
	Zip      string `gorm:"not null"`
	City     string `gorm:"not null"`
	Province string `gorm:"not null"`
}

type Order struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	TotalPrice        float64   `gorm:"not null;type:numeric(8,2)"`
	PaymentReceivedAt time.Time
	OrderStatus       string    `gorm:"not null"`
	PaymentStatus     string    `gorm:"not null"`
	UserID            uuid.UUID `gorm:"not null"`
	User              User      `gorm:"not null;foreignKey:UserID;references:ID"`
	UserName          string    `gorm:"not null"`
	UserEmail         string    `gorm:"not null"`
	UserPhoneNumber   string    `gorm:"not null"`
	Address           Address   `gorm:"embedded;embeddedPrefix:addr_"`
}

type OrderWithItems struct {
	Order
	OrderItems []*OrderItem
}
