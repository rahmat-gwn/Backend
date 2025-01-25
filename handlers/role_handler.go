package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func CheckRole(c *gin.Context) {
    userRole := c.GetHeader("X-User-Role")
    if userRole == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Role header missing"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Role retrieved successfully",
        "role":    userRole,
    })
}
