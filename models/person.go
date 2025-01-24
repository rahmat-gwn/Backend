package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address"`
    Phone       string `gorm:"unique" json:"phone"`
	IDNumber    string `json:"id_number" gorm:"unique;not null"`
	PhoneNumber string `json:"phone_number"`
}

// BeforeCreate hook untuk mengisi id_number secara otomatis
func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	if p.IDNumber == "" {
		p.IDNumber = uuid.New().String() // Menghasilkan UUID sebagai id_number
	}
	return
}
