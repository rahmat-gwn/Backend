package handlers

import (
    "go-modul/database"
    "go-modul/models"
    "net/http"
	"strings"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// RoleMiddleware ensures the user has the required role to access certain endpoints.
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

	// Bind JSON body ke struct person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID role dari database
	var userRole models.Role
	if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Default role not found"})
		return
	}

	// Set role ke user jika RoleID kosong
	if person.RoleID == 0 {
		person.RoleID = userRole.ID
	}

	// Simpan data ke database
	if err := database.DB.Create(&person).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number or ID number already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Person created successfully",
		"data": map[string]interface{}{
			"id":      person.ID,
			"name":    person.Name,
			"address": person.Address,
			"phone":   person.Phone,
			"role_id": person.RoleID,
		},
	})
}

func UpdatePerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person

	// Periksa apakah data person ada
	if err := database.DB.First(&person, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	// Ambil data dari request body
	var input models.Person
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID role berdasarkan nama jika diubah
	if input.RoleID != 0 {
		var role models.Role
		if err := database.DB.First(&role, input.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}
		person.RoleID = role.ID
	}

	// Update field lainnya
	person.Name = input.Name
	person.Address = input.Address
	person.Phone = input.Phone

	// Simpan perubahan ke database
	if err := database.DB.Save(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person updated successfully", "data": person})
}

func DeletePerson(c *gin.Context) {
    id := c.Param("id")
    var person models.Person

    if err := database.DB.First(&person, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Person not found or already deleted"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check person"})
        }
        return
    }

    if err := database.DB.Delete(&person).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}
