package routes

import (
	"go-modul/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua rute API
func RegisterRoutes(r *gin.Engine) {
	r.GET("/people", handlers.GetPeople)
	r.GET("/people/:id", handlers.GetPersonByID)
	r.POST("/people", handlers.CreatePerson)
	r.PUT("/people/:id", handlers.UpdatePerson)
	r.DELETE("/people/:id", handlers.DeletePerson)
}
