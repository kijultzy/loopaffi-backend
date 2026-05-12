package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// ==========================================
	// 🔴 RUTE PUBLIK (Bisa diakses tanpa token)
	// ==========================================
	api.POST("/auth/login", controllers.Login)
	api.POST("/roles", controllers.CreateRole)
	api.POST("/users", controllers.CreateUser) // Biasanya register user dibiarkan publik

	// ==========================================
	// 🟢 RUTE TERPROTEKSI (Wajib pakai token)
	// ==========================================
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Pasang satpam JWT di sini
	{
		// Users (Read-only)
		protected.GET("/users", controllers.GetUsers)
		protected.GET("/users/:id", controllers.GetUserByID)

		// Sales (Read-only)
		protected.GET("/sales", controllers.GetSales)
		protected.GET("/sales/:id", controllers.GetSaleByID)

		// Commissions (Read-only)
		protected.GET("/commissions", controllers.GetCommissions)
		protected.GET("/commissions/:id", controllers.GetCommissionByID)

		// Payments (Read-only)
		protected.GET("/payments", controllers.GetPayments)

		// Notifications & Dashboard
		protected.GET("/notifications", controllers.GetNotifications)
		protected.PUT("/notifications/:id/read", controllers.MarkNotificationRead)
		protected.GET("/dashboard/stats", controllers.GetDashboardStats)

		// ==========================================
		// 🔴 RUTE KHUSUS ADMIN (Wajib token + Admin)
		// ==========================================
		admin := protected.Group("/")
		admin.Use(middleware.AdminOnly()) // Satpam kedua, hanya untuk admin
		{
			// Mutasi Sales
			admin.POST("/sales", controllers.CreateSale)

			// Endpoint Komisi
			admin.POST("/commissions/calculate", controllers.HitungKomisi)

			// Endpoint Pembayaran
			admin.POST("/payments/process", controllers.ProsesPembayaran)
			admin.PUT("/payments/:id/pay", controllers.MarkPaymentPaid)

			// Laporan
			admin.GET("/reports", controllers.BuatLaporan)
		}
	}
}