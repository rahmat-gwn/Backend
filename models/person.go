package models

type Person struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Address     string `gorm:"size:255"`
	Phone       string `gorm:"size:50"`
	IDNumber    string `gorm:"size:50;uniqueIndex"`
	PhoneNumber string `gorm:"size:50"`
}
