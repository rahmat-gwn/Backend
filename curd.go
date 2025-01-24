package main

import (
	"go-modul/database"
    "go-modul/handlers"
    "go-modul/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi database
	database.InitDatabase()

	// Migrasi model ke database
	database.DB.AutoMigrate(&models.Person{})

	// Setup router
	r := gin.Default()

	// Rute API
	r.GET("/people", handlers.GetPeople)
	r.GET("/people/:id", handlers.GetPersonByID)
	r.POST("/people", handlers.CreatePerson)
	r.PUT("/people/:id", handlers.UpdatePerson)
	r.DELETE("/people/:id", handlers.DeletePerson)

	// Jalankan server sesuaikan jika 8080 sudah terpakai
	r.Run(":8081")
}
