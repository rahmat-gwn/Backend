package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"go-modul/models"
)

var DB *gorm.DB

// InitDatabase untuk menginisialisasi koneksi ke database
func InitDatabase() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Open database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	fmt.Println("Database connection established")
}

// Migrate untuk melakukan migrasi model tanpa membuat tabel yang sama
func Migrate(models ...interface{}) {
	for _, model := range models {
		// Pengecekan jika model sudah ada, jangan buat lagi
		if err := DB.AutoMigrate(model); err != nil {
			log.Printf("Failed to migrate model %v: %v", model, err)
		} else {
			fmt.Printf("Migrated model %v successfully\n", model)
		}
	}
}

func GetDBName() string {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is not set in .env file")
	}
	return dbName
}

func GetNextAvailableID() (int, error) {
    var ids []int
    if err := DB.Model(&models.Person{}).Pluck("id", &ids).Error; err != nil {
        return 0, err
    }

    for i := 1; ; i++ {
        if !contains(ids, i) {
            return i, nil
        }
    }
}

// Fungsi pembantu untuk memeriksa apakah ID ada dalam daftar
func contains(slice []int, item int) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}