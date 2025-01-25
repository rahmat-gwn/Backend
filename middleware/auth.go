package middleware

import (
    "go-modul/handlers"
    "go-modul/middleware"
    "github.com/gin-gonic/gin"
	"net/http"
)

//"net/http"
// func RegisterRoutes(r *gin.Engine) {
//     personRoutes := r.Group("/people")
//     {
//         // Role admin bisa melakukan semua tindakan
//         personRoutes.Use(middleware.RoleAuthorization("admin"))
//         personRoutes.GET("/", handlers.GetPeople)
//         personRoutes.PUT("/:id", handlers.UpdatePerson)
//         personRoutes.DELETE("/:id", handlers.DeletePerson)
//     }

//     userRoutes := r.Group("/people/user")
//     {
//         // Role user hanya bisa mengedit nama, alamat, dan phone
//         userRoutes.Use(middleware.RoleAuthorization("admin", "user"))
//         userRoutes.PUT("/:id", handlers.UpdatePersonForUser)
//     }
// }

func RoleAuthorization(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Ambil role dari header
        userRole := c.GetHeader("X-User-Role")
        if userRole == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this resource"})
            c.Abort()
            return
        }

        // Periksa apakah role diperbolehkan
        for _, role := range allowedRoles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this resource"})
        c.Abort()
    }
}

func RoleMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Contoh sederhana: Ambil role dari header
        role := c.GetHeader("X-User-Role")
        if role == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not provided"})
            c.Abort()
            return
        }

        // Simpan role ke context
        c.Set("role", role)
        c.Next()
    }
}
