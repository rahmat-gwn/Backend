package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Person struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Name        string `gorm:"not null" json:"name"`
    Address     string `json:"address"`
    Phone       string `gorm:"unique;not null" json:"phone"`
    IDNumber    string `gorm:"unique;not null" json:"id_number"`
    PhoneNumber string `json:"phone_number"`
    RoleID      uint   `gorm:"not null" json:"role_id"`
    Role        Role   `gorm:"foreignKey:RoleID" json:"role"`
}
// BeforeCreate hook untuk mengisi id_number secara otomatis
func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	if p.IDNumber == "" {
		p.IDNumber = uuid.New().String() // Menghasilkan UUID sebagai id_number
	}
	return
}
