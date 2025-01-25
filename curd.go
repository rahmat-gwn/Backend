package main

import (
	"go-modul/database"
	"go-modul/models"
	"go-modul/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi database
	database.InitDatabase()

	// Migrasi model ke database
	if err := database.DB.AutoMigrate(&models.Person{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Setup router
	r := gin.Default()

	// Registrasi rute melalui modul routes
	routes.RegisterRoutes(r)

	// Jalankan server (ubah port jika 8081 sudah terpakai)
	if err := r.Run(":8081"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
