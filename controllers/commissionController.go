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

// Menangkap request dari Admin
type CommissionInput struct {
	ID     string `json:"id_commission"`
	SaleID string `json:"id_sale"` // ID Penjualan mana yang mau dihitung komisinya
}

// ==================== Handlers ====================

// POST /api/commissions/calculate — Hitung komisi untuk penjualan tertentu
// Implementasi Use Case #2: Hitung Komisi
func HitungKomisi(c *gin.Context) {
	var input CommissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Tarik data Sale beserta isi SaleItems-nya
	var sale models.Sale
	if err := config.DB.Preload("SaleItems").First(&sale, "id_sale = ?", input.SaleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data penjualan tidak ditemukan!"})
		return
	}

	// Hitung total dari TotalAmount atau dari SaleItems
	var totalPenjualan float64
	if sale.TotalAmount > 0 {
		totalPenjualan = sale.TotalAmount
	} else {
		for _, item := range sale.SaleItems {
			// Gunakan keyword break jika ada anomali nilai negatif
			if item.Subtotal < 0 {
				break
			}
			totalPenjualan += item.Subtotal
		}
	}

	// 2. Tarik CommissionSetting yang is_active = true
	var setting models.CommissionSetting
	if err := config.DB.Where("is_active = ?", true).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada setting komisi yang aktif saat ini!"})
		return
	}

	// 3. Algoritma: amount = total * (persentase / 100)
	// Mengambil nilai dari persentase aktif (getPersentaseAktif())
	jumlahKomisi := totalPenjualan * (setting.PersentaseKomisi / 100)

	// 4. Generate ID kalau tidak disediakan
	commissionID := input.ID
	if commissionID == "" {
		commissionID = fmt.Sprintf("COM-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
	}

	// 5. Buat objek komisi baru
	commission := models.Commission{
		ID:                  commissionID,
		SaleID:              sale.ID,
		AffiliateID:         sale.UserID,
		CommissionSettingID: setting.ID,
		JumlahKomisi:        jumlahKomisi,
		TglHitung:           time.Now(),
		StatusKomisi:        "Pending",
	}

	// Simpan ke database
	if err := config.DB.Create(&commission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data komisi!"})
		return
	}

	// Berikan response sukses
	c.JSON(http.StatusOK, gin.H{
		"message": "Komisi berhasil dihitung dan berstatus Pending!",
		"data":    commission,
		"rincian_kalkulasi": gin.H{
			"total_penjualan":  totalPenjualan,
			"persentase_aktif": setting.PersentaseKomisi,
		},
	})
}

// GET /api/commissions — Ambil semua data komisi
func GetCommissions(c *gin.Context) {
	var commissions []models.Commission

	// Query parameter optional: affiliateId untuk filter per affiliate
	affiliateID := c.Query("affiliateId")

	query := config.DB.Order("tgl_hitung DESC")
	if affiliateID != "" {
		query = query.Where("id_affiliate = ?", affiliateID)
	}

	if err := query.Find(&commissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data komisi!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data komisi berhasil diambil!",
		"data":    commissions,
	})
}

// GET /api/commissions/:id — Ambil detail komisi by ID
func GetCommissionByID(c *gin.Context) {
	id := c.Param("id")

	var commission models.Commission
	if err := config.DB.Preload("Sale").Preload("CommissionSetting").First(&commission, "id_commission = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data komisi tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data komisi berhasil diambil!",
		"data":    commission,
	})
}