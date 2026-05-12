package middleware

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil id_user dari context yang diset oleh AuthMiddleware
		idUser, exists := c.Get("id_user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! Token tidak ditemukan atau tidak valid."})
			c.Abort()
			return
		}

		var user models.User
		// Cari user beserta Role-nya
		if err := config.DB.Preload("Role").First(&user, "id_user = ?", idUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak! User tidak ditemukan."})
			c.Abort()
			return
		}

		// Cek apakah Role adalah Admin
		if user.Role.NamaRole != "Admin" && user.Role.NamaRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses diblokir! Hanya admin yang diizinkan mengakses rute ini."})
			c.Abort()
			return
		}

		// Jika admin, lanjut ke handler berikutnya
		c.Next()
	}
}
