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

	// Cek apakah tabel sudah ada, jika tidak maka migrasi
	if !database.DB.Migrator().HasTable(&models.Role{}) {
		err := database.DB.AutoMigrate(&models.Role{})
		if err != nil {
			log.Fatalf("Failed to migrate table Role: %v", err)
		}
		fmt.Println("Table 'roles' migrated successfully.")

		// Seed data untuk tabel Role
		roles := []models.Role{
			{Name: "administrator"},
			{Name: "user"},
		}

		for _, role := range roles {
			if err := database.DB.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role %s: %v", role.Name, err)
			} else {
				fmt.Printf("Seeded role: %s\n", role.Name)
			}
		}
	} else {
		fmt.Println("Table 'roles' already exists, skipping migration.")
	}

	if !database.DB.Migrator().HasTable(&models.Person{}) {
		err := database.DB.AutoMigrate(&models.Person{})
		if err != nil {
			log.Fatalf("Failed to migrate table Person: %v", err)
		}
		fmt.Println("Table 'people' migrated successfully.")
	} else {
		fmt.Println("Table 'people' already exists, skipping migration.")
	}

	fmt.Println("Migration completed.")
}
