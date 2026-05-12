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

type PaymentInput struct {
	ID              string `json:"id_payment"`
	CommissionID    string `json:"id_commission"`
	PaymentMethodID string `json:"id_payment_method"`
}

// Implementasi Use Case #3: Kelola Pembayaran
func ProsesPembayaran(c *gin.Context) {
	var input PaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var commission models.Commission
	if err := config.DB.Preload("Sale").First(&commission, "id_commission = ?", input.CommissionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data komisi tidak ditemukan!"})
		return
	}

	if commission.StatusKomisi == "Lunas" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Komisi ini sudah dibayar!"})
		return
	}

	tx := config.DB.Begin()

	paymentID := input.ID
	if paymentID == "" {
		paymentID = fmt.Sprintf("PAY-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
	}

	payment := models.Payment{
		ID: paymentID, CommissionID: commission.ID, AffiliateID: commission.AffiliateID,
		PaymentMethodID: input.PaymentMethodID, JumlahBayar: commission.JumlahKomisi,
		TglPembayaran: time.Now(), StatusBayar: "Lunas",
	}
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pembayaran!"})
		return
	}

	commission.StatusKomisi = "Lunas"
	if err := tx.Save(&commission).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update komisi!"})
		return
	}

	affID := commission.AffiliateID
	if affID == "" {
		affID = commission.Sale.UserID
	}
	// Trigger pembuatan data Notifikasi (representasi dari fungsi sendNotif())
	notif := models.Notification{
		ID: fmt.Sprintf("NOTIF-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000)),
		UserID: affID, Judul: "Pembayaran Komisi Berhasil",
		Pesan: fmt.Sprintf("Komisi sebesar Rp %.0f telah dicairkan.", commission.JumlahKomisi),
		IsRead: false,
	}
	tx.Create(&notif)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil!", "data": payment})
}

func GetPayments(c *gin.Context) {
	var payments []models.Payment
	affiliateID := c.Query("affiliateId")
	query := config.DB.Order("tgl_pembayaran DESC")
	if affiliateID != "" {
		query = query.Where("id_affiliate = ?", affiliateID)
	}
	if err := query.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": payments})
}

func MarkPaymentPaid(c *gin.Context) {
	id := c.Param("id")
	var payment models.Payment
	if err := config.DB.First(&payment, "id_payment = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pembayaran tidak ditemukan!"})
		return
	}
	if payment.StatusBayar == "Lunas" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sudah lunas!"})
		return
	}

	tx := config.DB.Begin()
	payment.StatusBayar = "Lunas"
	payment.TglPembayaran = time.Now()
	tx.Save(&payment)

	if payment.CommissionID != "" {
		tx.Model(&models.Commission{}).Where("id_commission = ?", payment.CommissionID).Update("status_komisi", "Lunas")
	}
	if payment.AffiliateID != "" {
		n := models.Notification{
			ID: fmt.Sprintf("NOTIF-%s-%04d", time.Now().Format("20060102150405"), rand.Intn(10000)),
			UserID: payment.AffiliateID, Judul: "Pembayaran Lunas",
			Pesan: fmt.Sprintf("Pembayaran Rp %.0f lunas.", payment.JumlahBayar), IsRead: false,
		}
		tx.Create(&n)
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran ditandai lunas!", "data": payment})
}