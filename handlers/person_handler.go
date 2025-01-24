package handlers

import (
    "go-modul/database"
    "go-modul/models"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// Middleware to check role permissions
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role") // Assume role is stored in context
        if !exists || (requiredRole != "any" && userRole != requiredRole) {
            c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func GetPeople(c *gin.Context) {
    var people []models.Person
    if err := database.DB.Find(&people).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
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

    if err := c.ShouldBindJSON(&person); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if phone number already exists
    var existingPerson models.Person
    if err := database.DB.Where("phone = ?", person.Phone).First(&existingPerson).Error; err == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number already exists"})
        return
    }

    // Set default role if not provided
    if person.RoleID == 0 {
        var userRole models.Role
        if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign default role"})
            return
        }
        person.RoleID = userRole.ID
    } else {
        // Validate provided role
        var role models.Role
        if err := database.DB.First(&role, person.RoleID).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
            return
        }
    }

    // Auto-generate ID number
    person.IDNumber = ""

    if err := database.DB.Create(&person).Error; err != nil {
        if strings.Contains(err.Error(), "duplicate") {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number or ID number already exists"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
        }
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Person created successfully", "data": person})
}

func UpdatePerson(c *gin.Context) {
    id := c.Param("id")
    var person models.Person

    if err := database.DB.First(&person, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
        return
    }

    var input models.Person
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if phone number is already used by another person
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
