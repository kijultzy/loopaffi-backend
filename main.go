package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Koneksi Database PostgreSQL
	config.ConnectDatabase()

	// 2. Auto-Migrate semua tabel
	err := config.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Product{},
		&models.PaymentMethod{},
		&models.CommissionSetting{},
		&models.Sale{},
		&models.SaleItem{},
		&models.Commission{},
		&models.Payment{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}
	fmt.Println("✅ Migrasi 10 tabel berhasil!")

	// 3. Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 4. Inisialisasi router
	router := gin.Default()

	// 5. CORS Middleware — Production-safe
	// Baca origin yang diizinkan dari env variable.
	// Di production (Render), set ALLOWED_ORIGIN=https://your-app.vercel.app
	// Di development, biarkan kosong untuk mengizinkan semua origin.
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowedOrigin != "" {
			// PRODUCTION: hanya izinkan domain frontend yang terdaftar
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			// DEVELOPMENT: izinkan semua origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Max-Age", "43200")

		// Tangani preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 6. Daftarkan routes
	routes.SetupRoutes(router)

	// 7. Jalankan server — Port dinamis untuk Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("🚀 LoopAffi backend berjalan di http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
