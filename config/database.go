package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// PERUBAHAN 1: Jangan abaikan error sepenuhnya.
	// Kita log error-nya agar tahu jika .env gagal dibaca saat mode development lokal.
	err := godotenv.Load()
	if err != nil {
		log.Println("INFO: File .env tidak ditemukan, mencoba menggunakan environment variables sistem (Production Mode).")
	}

	// Prioritaskan DATABASE_URL (format untuk deploy di Railway/Render nanti)
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		// Fallback ke variabel terpisah
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		// PERUBAHAN 2: Tambahkan validasi. Hentikan program jika variabel penting kosong.
		// Ini mencegah error "invalid port" yang membingungkan dari GORM.
		if host == "" || port == "" || dbname == "" {
			log.Fatal("ERROR: Konfigurasi database tidak lengkap! Pastikan DB_HOST, DB_PORT, dan DB_NAME ada di file .env")
		}

		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			host, user, password, dbname, port,
		)
	}

	// Buka koneksi pakai GORM
	database, gormErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if gormErr != nil {
		log.Fatal("Gagal koneksi ke database PostgreSQL! Detail error: ", gormErr)
	}

	fmt.Println("Koneksi database PostgreSQL berhasil!")
	DB = database
}