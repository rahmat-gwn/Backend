package main

import (
	"fmt"
	"go-modul/database"
	"go-modul/models"
	"github.com/google/uuid"
)

func main() {
	// Inisialisasi koneksi database
	database.InitDatabase()

	// Seed data untuk tabel Person
	people := []models.Person{
		{Name: "John Doe", Address: "123 Elm Street", Phone: "123-456-7890", IDNumber: uuid.New().String()},
		{Name: "Jane Smith", Address: "456 Oak Avenue", Phone: "987-654-3210", IDNumber: uuid.New().String()},
		{Name: "Alice Brown", Address: "789 Pine Lane", Phone: "555-123-4567", IDNumber: uuid.New().String()},
		{Name: "Bob Johnson", Address: "321 Maple Drive", Phone: "444-555-6666", IDNumber: uuid.New().String()},
		{Name: "Charlie Davis", Address: "654 Cedar Court", Phone: "333-777-8888", IDNumber: uuid.New().String()},
	}

	// Tambahkan data ke tabel Person
	for _, person := range people {
		var existingPerson models.Person
		// Pastikan IDNumber unik
		if err := database.DB.Where("id_number = ?", person.IDNumber).FirstOrCreate(&existingPerson, person).Error; err != nil {
			fmt.Printf("Failed to insert person %s: %v\n", person.Name, err)
		} else {
			fmt.Printf("Inserted/Checked person: %s\n", person.Name)
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
