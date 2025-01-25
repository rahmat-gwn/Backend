package database

import (
	"fmt"
	"go-modul/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	// Log connection success
	fmt.Println("Database connection established")

	// Migrasi tabel
	err = DB.AutoMigrate(&models.Role{}, &models.Person{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Tambahkan data awal untuk tabel roles
	seedRoles()
}

// seedRoles menambahkan role default jika belum ada.
func seedRoles() {
	var count int64
	DB.Model(&models.Role{}).Count(&count)

	if count == 0 {
		roles := []models.Role{
			{Name: "administrator"},
			{Name: "user"},
		}

		if err := DB.Create(&roles).Error; err != nil {
			log.Fatalf("Failed to seed roles: %v", err)
		}

		fmt.Println("Default roles seeded successfully")
	}
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

func GetNextAvailableID() (uint, error) {
    var nextID uint

    // Query for the smallest missing ID
    query := `
        SELECT MIN(t1.id + 1) AS next_id
        FROM people t1
        LEFT JOIN people t2 ON t1.id + 1 = t2.id
        WHERE t2.id IS NULL
    `
    row := DB.Raw(query).Row()
    if err := row.Scan(&nextID); err != nil {
        return 0, err
    }

    // If no gaps, return the next sequential ID
    if nextID == 0 {
        var maxID uint
        DB.Table("people").Select("MAX(id)").Row().Scan(&maxID)
        nextID = maxID + 1
    }

    return nextID, nil
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