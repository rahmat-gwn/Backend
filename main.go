package main

import (
	"fmt"
	"log"
	"go-modul/database"
	"go-modul/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	database.InitDatabase()

	// Periksa apakah database sudah ada
	var count int64
	err = database.DB.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", database.GetDBName()).Scan(&count).Error
	if err != nil || count == 0 {
		log.Println("Database tidak ditemukan. Membuat database...")

		// Buat database jika belum ada
		err := database.DB.Exec("CREATE DATABASE IF NOT EXISTS ?", database.GetDBName()).Error
		if err != nil {
			log.Fatalf("Error membuat database: %v", err)
		}
	} else {
		log.Println("Database sudah ada, melanjutkan migrasi...")
	}

	// Periksa apakah tabel sudah ada
	var hasTable bool
	err = database.DB.Raw("SHOW TABLES LIKE 'people'").Scan(&hasTable).Error
	if err != nil || !hasTable {
		log.Println("Tabel tidak ditemukan. Melakukan migrasi...")

		// Lakukan AutoMigrate untuk membuat tabel jika belum ada
		err = database.DB.AutoMigrate(
			&models.Person{},
			&models.Product{},
		)
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		log.Println("Migrasi berhasil dilakukan!")
	} else {
		log.Println("Tabel sudah ada, melanjutkan seeding...")
	}

	// Seed data
	seedData()
}

func seedData() {
	// Seed data untuk tabel Person
	people := []models.Person{
		{Name: "John Doe", Address: "123 Elm Street", Phone: "123-456-7890"},
		{Name: "Jane Smith", Address: "456 Oak Avenue", Phone: "987-654-3210"},
		{Name: "Alice Brown", Address: "789 Pine Lane", Phone: "555-123-4567"},
		{Name: "Bob Johnson", Address: "321 Maple Drive", Phone: "444-555-6666"},
		{Name: "Charlie Davis", Address: "654 Cedar Court", Phone: "333-777-8888"},
	}

	// Tambahkan data ke tabel Person
	for _, person := range people {
		var existingPerson models.Person
		err := database.DB.Where("id_number = ?", person.IDNumber).First(&existingPerson).Error
		if err != nil { // Data tidak ditemukan, insert data baru
			if err := database.DB.Create(&person).Error; err != nil {
				fmt.Printf("Failed to insert person %s: %v\n", person.Name, err)
			} else {
				fmt.Printf("Inserted person: %s\n", person.Name)
			}
		} else {
			// Jika sudah ada, tidak lakukan apa-apa
			fmt.Printf("Person %s already exists, skipping.\n", person.Name)
		}
	}

	// Seed data untuk tabel Product
	products := []models.Product{
		{Name: "Laptop", Price: 1200.00},
		{Name: "Smartphone", Price: 800.00},
		{Name: "Headphones", Price: 150.00},
		{Name: "Monitor", Price: 300.00},
		{Name: "Keyboard", Price: 50.00},
	}

	// Tambahkan data ke tabel Product
	for _, product := range products {
		if err := database.DB.Create(&product).Error; err != nil {
			fmt.Printf("Failed to insert product %s: %v\n", product.Name, err)
		} else {
			fmt.Printf("Inserted product: %s\n", product.Name)
		}
	}
}
