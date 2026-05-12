package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LaporanResult struct {
	NamaAffiliate  string  `json:"nama_affiliate"`
	Email          string  `json:"email"`
	TotalTransaksi int     `json:"total_transaksi"`
	TotalPenjualan float64 `json:"total_penjualan"`
	TotalKomisi    float64 `json:"total_komisi"`
	KomisiDibayar  float64 `json:"komisi_dibayar"`
	KomisiTertunda float64 `json:"komisi_tertunda"`
}

// GET /api/reports
// Implementasi Use Case #4: Buat Laporan (Ekstraksi rekapitulasi performa seluruh afiliasi)
func BuatLaporan(c *gin.Context) {
	var hasilLaporan []LaporanResult

	query := `
		SELECT 
			u.nama_user as nama_affiliate,
			u.email as email,
			COUNT(DISTINCT s.id_sale) as total_transaksi, 
			COALESCE(SUM(s.total_amount), 0) as total_penjualan, 
			COALESCE(SUM(c.jumlah_komisi), 0) as total_komisi,
			COALESCE(SUM(CASE WHEN p.status_bayar = 'Lunas' THEN p.jumlah_bayar ELSE 0 END), 0) as komisi_dibayar,
			COALESCE(SUM(CASE WHEN p.status_bayar = 'pending' THEN p.jumlah_bayar ELSE 0 END), 0) as komisi_tertunda
		FROM users u
		INNER JOIN roles r ON u.id_role = r.id_role
		LEFT JOIN sales s ON u.id_user = s.id_user
		LEFT JOIN commissions c ON s.id_sale = c.id_sale
		LEFT JOIN payments p ON c.id_commission = p.id_commission
		WHERE r.nama_role != 'Admin' AND r.nama_role != 'admin'
		GROUP BY u.id_user, u.nama_user, u.email
	`

	if err := config.DB.Raw(query).Scan(&hasilLaporan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik data laporan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Laporan performa Affiliate berhasil ditarik!",
		"data":    hasilLaporan,
	})
}

// GET /api/notifications
func GetNotifications(c *gin.Context) {
	var notifications []models.Notification
	userID := c.Query("userId")

	query := config.DB.Order("created_at DESC")
	if userID != "" {
		query = query.Where("id_user = ?", userID)
	}

	if err := query.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil notifikasi!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": notifications})
}

// PUT /api/notifications/:id/read
func MarkNotificationRead(c *gin.Context) {
	id := c.Param("id")

	var notif models.Notification
	if err := config.DB.First(&notif, "id_notification = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notifikasi tidak ditemukan!"})
		return
	}

	notif.IsRead = true
	config.DB.Save(&notif)

	c.JSON(http.StatusOK, gin.H{"message": "Notifikasi ditandai sudah dibaca!", "data": notif})
}

// GET /api/dashboard/stats
func GetDashboardStats(c *gin.Context) {
	affiliateID := c.Query("affiliateId")

	var totalSales float64
	var salesCount int64
	var totalCommissions float64
	var pendingPayments float64
	var paidPayments float64
	var pendingCount int64
	var paidCount int64

	// Query Sales
	salesQ := config.DB.Model(&models.Sale{})
	if affiliateID != "" {
		salesQ = salesQ.Where("id_user = ?", affiliateID)
	}
	salesQ.Count(&salesCount)
	salesQ.Select("COALESCE(SUM(total_amount), 0)").Scan(&totalSales)

	// Query Commissions
	commQ := config.DB.Model(&models.Commission{})
	if affiliateID != "" {
		commQ = commQ.Where("id_affiliate = ?", affiliateID)
	}
	commQ.Select("COALESCE(SUM(jumlah_komisi), 0)").Scan(&totalCommissions)

	// Query Pending Payments
	pendingPayQ := config.DB.Model(&models.Payment{}).Where("status_bayar = ?", "pending")
	if affiliateID != "" {
		pendingPayQ = pendingPayQ.Where("id_affiliate = ?", affiliateID)
	}
	pendingPayQ.Count(&pendingCount)
	pendingPayQ.Select("COALESCE(SUM(jumlah_bayar), 0)").Scan(&pendingPayments)

	// Query Paid Payments
	paidPayQ := config.DB.Model(&models.Payment{}).Where("status_bayar = ?", "Lunas")
	if affiliateID != "" {
		paidPayQ = paidPayQ.Where("id_affiliate = ?", affiliateID)
	}
	paidPayQ.Count(&paidCount)
	paidPayQ.Select("COALESCE(SUM(jumlah_bayar), 0)").Scan(&paidPayments)

	c.JSON(http.StatusOK, gin.H{
		"totalSalesAmount":    totalSales,
		"salesCount":          salesCount,
		"totalCommissions":    totalCommissions,
		"pendingPayments":     pendingPayments,
		"paidPayments":        paidPayments,
		"pendingPaymentCount": pendingCount,
		"paidPaymentCount":    paidCount,
	})
}