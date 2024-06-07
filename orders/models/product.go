package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name    string    `gorm:"not null"`
	Price   float64   `gorm:"type:numeric(8,2);not null"`
	InStock int       `gorm:"not null"`
}
