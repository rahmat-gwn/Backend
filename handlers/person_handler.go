package handlers

import (
    "go-modul/database"
    "go-modul/models"
    "net/http"
    "strings"
	"gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func GetPeople(c *gin.Context) {
    var people []models.Person
    if err := database.DB.Find(&people).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
        return
    }
    c.JSON(http.StatusOK, people)
}


func GetPersonByID(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	if err := database.DB.First(&person, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(http.StatusOK, person)
}

func CreatePerson(c *gin.Context) {
    var person models.Person

    // Bind the JSON body to the person struct
    if err := c.ShouldBindJSON(&person); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if phone number already exists in the database
    var existingPerson models.Person
    if err := database.DB.Where("phone = ?", person.Phone).First(&existingPerson).Error; err == nil {
        // If a person with the same phone exists, return an error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
        return
    }

    // Tentukan ID secara manual
    nextID, err := database.GetNextAvailableID()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to determine next ID"})
        return
    }
    person.ID = nextID

    // Bersihkan nilai id_number untuk diisi otomatis
    person.IDNumber = ""

    // Create the person in the database
    if err := database.DB.Create(&person).Error; err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number or ID number already exists"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
        }
        return
    }

    // Return success message
    c.JSON(http.StatusCreated, gin.H{"message": "Person created successfully", "data": person})
}

func UpdatePerson(c *gin.Context) {
    var person models.Person
    id := c.Param("id")

    // Check if the person exists
    if err := database.DB.First(&person, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
        return
    }

    var input models.Person
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if the phone number is already in use by someone else
    var existingPerson models.Person
    result := database.DB.Where("phone = ? AND id != ?", input.Phone, id).First(&existingPerson)
    if result.Error == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
        return
    } else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check phone number"})
        return
    }

    // Update person details
    person.Name = input.Name
    person.Address = input.Address
    person.Phone = input.Phone

    if err := database.DB.Save(&person).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully", "data": person})
}

func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Person{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Person deleted"})
}
