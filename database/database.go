package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk koneksi database
var DB *gorm.DB

// InitDatabase menginisialisasi koneksi ke database
func InitDatabase() {
	// Muat konfigurasi dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading environment variables from system")
	}

	// Format DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Hubungkan ke database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Assign koneksi ke variabel global
	DB = db
	fmt.Println("Database connection established")
}

// Migrate menjalankan migrasi database
func Migrate(models ...interface{}) {
	if DB == nil {
		log.Fatalf("Database not initialized. Call InitDatabase() first.")
	}

	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database migration completed successfully")
}
