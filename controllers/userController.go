package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ==================== LOGIN ====================

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// POST /api/auth/login
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input tidak valid: " + err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Pesan generik untuk mencegah user enumeration attack
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah!"})
		return
	}

	// Validasi password: coba bcrypt terlebih dahulu, fallback ke plain-text (kompatibilitas seed)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		// Fallback plain-text untuk data seed lama
		if user.PasswordHash != input.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah!"})
			return
		}
	}

	roleName := user.GetRoleName()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "loopaffi_rahasia_super_aman_ganti_di_production"
	}

	// Generate JWT dengan claims lengkap (termasuk role untuk RBAC di frontend)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_user": user.ID,
		"email":   user.Email,
		"role":    roleName,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token autentikasi!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   tokenString,
		"user":    user.ToResponse(),
	})
}

// ==================== ROLE ====================

// POST /api/roles
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat role!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Role berhasil dibuat!", "data": role})
}

// ==================== USER (REGISTER) ====================

type RegisterInput struct {
	ID           string `json:"id"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password_hash" binding:"required,min=6"`
	Phone        string `json:"phone"`
}

// POST /api/users — Registrasi publik (tidak butuh token)
func CreateUser(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Cek duplikat email sebelum insert
	var existing models.User
	if config.DB.Where("email = ?", input.Email).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar! Silakan gunakan email lain."})
		return
	}

	// Hash password dengan bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password!"})
		return
	}

	// Generate ID jika tidak disediakan client
	userID := input.ID
	if userID == "" {
		userID = "USR-" + time.Now().Format("20060102150405")
	}

	// Role SELALU Affiliate (ROLE-002) — admin tidak bisa dibuat lewat registrasi publik
	user := models.User{
		ID:           userID,
		RoleID:       "ROLE-002",
		NamaUser:     input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedBytes),
		NoHp:         input.Phone,
		StatusUser:   "active",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan akun! Coba lagi."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil! Silakan login.",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.NamaUser,
			"email": user.Email,
		},
	})
}

// GET /api/users
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data users!"})
		return
	}

	var responses []models.UserResponse
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data users berhasil diambil!", "data": responses})
}

// GET /api/users/:id
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.Preload("Role").First(&user, "id_user = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data user berhasil diambil!", "data": user.ToResponse()})
}