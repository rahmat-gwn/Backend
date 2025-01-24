package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Struct: Data Bentuk
type Bentuk struct {
	ID       string  `json:"id"`
	Jenis    string  `json:"jenis"`
	Parameter float64 `json:"parameter"`
	Luas     float64 `json:"luas"`
}

// Slice untuk menyimpan data
var dataBentuk []Bentuk

// Fungsi untuk menghitung luas berdasarkan jenis bentuk
func hitungLuas(jenis string, parameter float64) (float64, error) {
	switch jenis {
	case "persegi":
		return parameter * parameter, nil
	case "lingkaran":
		return math.Pi * parameter * parameter, nil
	default:
		return 0, fmt.Errorf("jenis bentuk tidak dikenali")
	}
}

func main() {
	r := gin.Default()

	// Endpoint untuk mendapatkan semua data
	r.GET("/bentuk", func(c *gin.Context) {
		c.JSON(http.StatusOK, dataBentuk)
	})

	// Endpoint untuk menambahkan data baru
	r.POST("/bentuk", func(c *gin.Context) {
		var input struct {
			Jenis     string  `json:"jenis"`
			Parameter float64 `json:"parameter"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
			return
		}

		// Hitung luas dan tambahkan data baru
		luas, err := hitungLuas(input.Jenis, input.Parameter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bentukBaru := Bentuk{
			ID:       uuid.New().String(),
			Jenis:    input.Jenis,
			Parameter: input.Parameter,
			Luas:     luas,
		}
		dataBentuk = append(dataBentuk, bentukBaru)

		c.JSON(http.StatusCreated, bentukBaru)
	})

	// Endpoint untuk menghapus data berdasarkan ID
	r.DELETE("/bentuk/:id", func(c *gin.Context) {
		id := c.Param("id")
		for i, b := range dataBentuk {
			if b.ID == id {
				dataBentuk = append(dataBentuk[:i], dataBentuk[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
	})

	// Jalankan server
	r.Run(":8080") // Server berjalan di localhost:8080
}
