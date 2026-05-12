package models

import "time"

// Tabel roles
type Role struct {
	ID       string `gorm:"primaryKey;type:varchar(50);column:id_role" json:"id"`
	NamaRole string `gorm:"type:varchar(50);column:nama_role" json:"name"`
}

// Tabel users
type User struct {
	ID           string    `gorm:"primaryKey;type:varchar(50);column:id_user" json:"id"`
	RoleID       string    `gorm:"type:varchar(50);column:id_role" json:"roleId"`
	NamaUser     string    `gorm:"type:varchar(100);column:nama_user" json:"name"`
	Email        string    `gorm:"type:varchar(100);unique;column:email" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);column:password_hash" json:"-"` // TIDAK dikirim ke client
	NoHp         string    `gorm:"type:varchar(20);column:no_hp" json:"phone"`
	StatusUser   string    `gorm:"type:varchar(50);column:status_user" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`

	// Relasi ke Role
	Role Role `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
}

// DTO — hanya field aman yang dikirim ke frontend
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"` // "admin" atau "affiliate" (lowercase)
}

// GetRoleName mengembalikan nama role lowercase yang konsisten
func (u *User) GetRoleName() string {
	if u.Role.NamaRole == "Admin" || u.Role.NamaRole == "admin" {
		return "admin"
	}
	return "affiliate"
}

// ToResponse mengkonversi User ke UserResponse (aman untuk dikirim ke client)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.NamaUser,
		Email: u.Email,
		Role:  u.GetRoleName(),
	}
}