package main

import (
	"fmt"
	"log"

	"go-modul/database"
	"go-modul/models"
)

func main() {
	// Inisialisasi koneksi ke database
	database.InitDatabase()

	// Migrasi tabel Role
	if err := database.DB.AutoMigrate(&models.Role{}); err != nil {
		log.Fatalf("Failed to migrate table Role: %v", err)
	}
	fmt.Println("Table 'roles' migrated successfully.")

	// Seed data untuk tabel Role jika belum ada
	roles := []models.Role{
		{Name: "administrator"},
		{Name: "user"},
	}
	for _, role := range roles {
		var existing models.Role
		if err := database.DB.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err := database.DB.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role %s: %v", role.Name, err)
			} else {
				fmt.Printf("Seeded role: %s\n", role.Name)
			}
		} else {
			fmt.Printf("Role %s already exists, skipping.\n", role.Name)
		}
	}

	// Migrasi tabel Person
	if err := database.DB.AutoMigrate(&models.Person{}); err != nil {
		log.Fatalf("Failed to migrate table Person: %v", err)
	}
	fmt.Println("Table 'people' migrated successfully.")

	fmt.Println("Migration completed.")
}
