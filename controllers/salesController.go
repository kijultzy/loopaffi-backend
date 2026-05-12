package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== DTOs ====================

// Struct khusus (DTO) untuk menangkap data JSON yang dikirim Frontend
// Sesuai format frontend: { date, amount, affiliateId, status }
type SaleInput struct {
	Date        string  `json:"date"`
	Amount      float64 `json:"amount" binding:"required"`
	AffiliateID string  `json:"affiliateId" binding:"required"`
	Status      string  `json:"status"`
}

// ==================== Handlers ====================

// POST /api/sales — Catat penjualan baru (sesuai frontend admin/sales)
// Implementasi Use Case #1: Input Data Penjualan (Proses pembuatan objek Sale dan hitungTotal())
func CreateSale(c *gin.Context) {
	var input SaleInput

	// 1. Tangkap data JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Parse tanggal, default ke sekarang jika kosong
	tglPenjualan := time.Now()
	if input.Date != "" {
		parsed, err := time.Parse(time.RFC3339, input.Date)
		if err == nil {
			tglPenjualan = parsed
		}
	}

	// 3. Default status
	status := input.Status
	if status == "" {
		status = "completed"
	}

	// 4. Generate ID otomatis
	saleID := fmt.Sprintf("SALE-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000))

	// 5. Buat objek Sale sesuai Model Entity
	sale := models.Sale{
		ID:           saleID,
		UserID:       input.AffiliateID,
		TglPenjualan: tglPenjualan,
		TotalAmount:  input.Amount,
		StatusSale:   status,
	}

	// 6. Simpan Data Penjualan
	if err := config.DB.Create(&sale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data penjualan!"})
		return
	}

	// 7. Otomatis hitung dan buat komisi
	var setting models.CommissionSetting
	if err := config.DB.Where("is_active = ?", true).First(&setting).Error; err == nil {
		// Hitung komisi
		jumlahKomisi := sale.TotalAmount * (setting.PersentaseKomisi / 100)
		commissionID := fmt.Sprintf("COM-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000))

		commission := models.Commission{
			ID:                  commissionID,
			SaleID:              sale.ID,
			AffiliateID:         sale.UserID,
			CommissionSettingID: setting.ID,
			JumlahKomisi:        jumlahKomisi,
			TglHitung:           time.Now(),
			StatusKomisi:        "Pending",
		}
		config.DB.Create(&commission)

		// Buat payment record (pending)
		paymentID := fmt.Sprintf("PAY-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
		payment := models.Payment{
			ID:            paymentID,
			CommissionID:  commissionID,
			AffiliateID:   sale.UserID,
			JumlahBayar:   jumlahKomisi,
			TglPembayaran: time.Now(),
			StatusBayar:   "pending",
		}
		config.DB.Create(&payment)

		// Buat notifikasi untuk affiliate
		notifID := fmt.Sprintf("NOTIF-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
		notification := models.Notification{
			ID:     notifID,
			UserID: sale.UserID,
			Judul:  "Penjualan Baru Tercatat",
			Pesan:  fmt.Sprintf("New sale recorded! You earned Rp %s", formatRupiah(jumlahKomisi)),
			IsRead: false,
		}
		config.DB.Create(&notification)
	}

	// 8. Kembalikan Response dalam format yang sama dengan frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Data penjualan berhasil diinput!",
		"data":    sale,
	})
}

// GET /api/sales — Ambil semua data penjualan
func GetSales(c *gin.Context) {
	var sales []models.Sale

	// Query parameter optional: affiliateId untuk filter per affiliate
	affiliateID := c.Query("affiliateId")

	query := config.DB.Order("tgl_penjualan DESC")
	if affiliateID != "" {
		query = query.Where("id_user = ?", affiliateID)
	}

	if err := query.Find(&sales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penjualan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penjualan berhasil diambil!",
		"data":    sales,
	})
}

// GET /api/sales/:id — Ambil detail penjualan by ID
func GetSaleByID(c *gin.Context) {
	id := c.Param("id")

	var sale models.Sale
	if err := config.DB.Preload("SaleItems").Preload("User").First(&sale, "id_sale = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data penjualan tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penjualan berhasil diambil!",
		"data":    sale,
	})
}

// Helper format rupiah untuk notifikasi
func formatRupiah(amount float64) string {
	return fmt.Sprintf("%.0f", amount)
}