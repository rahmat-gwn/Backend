package main

import (
	"go-modul/database"
	"go-modul/models"
)

func main() {
	// Inisialisasi koneksi database
	database.InitDatabase()

	// Lakukan migrasi untuk semua model yang dibutuhkan
	database.Migrate(
		&models.Person{},
		&models.Product{}, // Model tambahan
	)
}
