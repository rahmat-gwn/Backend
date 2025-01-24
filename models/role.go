package models

type Role struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `gorm:"unique;not null" json:"name"`
}
