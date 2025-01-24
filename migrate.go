package main

import (
	"go-modul/database"
	"go-modul/models"
)

func main() {
	// Inisialisasi database
	database.InitDatabase()

	// Migrasi tabel
	database.Migrate(
		&models.Person{},
		&models.Product{},
	)
}
