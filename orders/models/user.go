package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"unique;not null"`
	PhoneNumber string    `gorm:"unique;not null"`
}
