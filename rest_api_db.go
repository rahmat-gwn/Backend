package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files" // Alias yang benar
	"github.com/swaggo/gin-swagger"       // Swagger UI handler
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Struktur untuk tabel Bentuk
type Bentuk struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	Jenis     string  `json:"jenis"`
	Parameter float64 `json:"parameter"`
	Luas      float64 `json:"luas"`
}

// Struktur untuk input data Bentuk
type BentukInput struct {
	Jenis     string  `json:"jenis" validate:"required,oneof=persegi lingkaran"`
	Parameter float64 `json:"parameter" validate:"required,gt=0"`
}

var validate = validator.New()

// Fungsi untuk menghitung luas
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

// Rate limiter sederhana
type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

type visitor struct {
	limiter  *time.Ticker
	lastSeen time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
	}
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := time.NewTicker(time.Minute / 200) 
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return rl.visitors[ip]
	}

	v.lastSeen = time.Now()
	return v
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > time.Minute {
				v.limiter.Stop()
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func rateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	go rl.cleanupVisitors()
	return func(c *gin.Context) {
		// Pengecualian untuk Swagger UI dan dokumentasi
		if c.Request.URL.Path == "/swagger/index.html" || 
			c.Request.URL.Path == "/swagger/swagger-ui.css" || 
			c.Request.URL.Path == "/swagger/swagger-ui-bundle.js" || 
			c.Request.URL.Path == "/swagger/swagger-ui-standalone-preset.js" || 
			c.Request.URL.Path == "/swagger/favicon-32x32.png" || 
			c.Request.URL.Path == "/swagger/favicon-16x16.png" || 
			c.Request.URL.Path == "/swagger/doc.json" || 
			c.Request.URL.Path == "/swagger/*any" {
			c.Next()
			return
		}

		ip := c.ClientIP()
		visitor := rl.getVisitor(ip)

		select {
		case <-visitor.limiter.C:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Terlalu banyak permintaan, coba lagi nanti"})
		}
	}
}



// @title API Bentuk
// @version 1.0
// @description API untuk menghitung dan menyimpan data bentuk geometris
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Memuat file .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Gagal memuat file .env: %v", err)
	}

	// Koneksi ke MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Migrasi otomatis
	db.AutoMigrate(&Bentuk{})

	// Set up Gin router
	r := gin.Default()

	// Inisialisasi rate limiter
	rl := NewRateLimiter()
	r.Use(rateLimitMiddleware(rl))

	// Middleware untuk logging request
	r.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		status := c.Writer.Status()
		log.Printf("Response: %d %s", status, http.StatusText(status))
	})

	// Menambahkan endpoint untuk dokumentasi Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Endpoint: Mendapatkan semua data bentuk
	r.GET("/bentuk", func(c *gin.Context) {
		var bentuk []Bentuk
		db.Find(&bentuk)
		c.JSON(http.StatusOK, bentuk)
	})

	// Endpoint: Menambahkan data baru
	r.POST("/bentuk", func(c *gin.Context) {
		var input BentukInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
			return
		}

		// Validasi input
		if err := validate.Struct(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hitung luas
		luas, err := hitungLuas(input.Jenis, input.Parameter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simpan data ke database
		bentukBaru := Bentuk{
			Jenis:     input.Jenis,
			Parameter: input.Parameter,
			Luas:      luas,
		}
		db.Create(&bentukBaru)

		c.JSON(http.StatusCreated, bentukBaru)
	})

	// Endpoint: Menghapus data berdasarkan ID
	r.DELETE("/bentuk/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Bentuk{}, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
	})

	// Menjalankan server di port 8080
	r.Run(":8080")
}
