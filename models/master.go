package models

import "time"

// Model untuk tabel products
type Product struct {
	ID            string  `gorm:"primaryKey;type:varchar(50);column:id_product" json:"id"`
	NamaProduct   string  `gorm:"type:varchar(100);column:nama_product" json:"name"`
	SKU           string  `gorm:"type:varchar(50);column:sku" json:"sku"`
	HargaDefault  float64 `gorm:"type:decimal(15,2);column:harga_default" json:"price"`
	StatusProduct string  `gorm:"type:varchar(50);column:status_product" json:"status"`
}

// Model untuk tabel payment_methods
type PaymentMethod struct {
	ID         string `gorm:"primaryKey;type:varchar(50);column:id_payment_method" json:"id"`
	NamaMetode string `gorm:"type:varchar(100);column:nama_metode" json:"name"`
}

// Model untuk tabel commission_settings
type CommissionSetting struct {
	ID               string    `gorm:"primaryKey;type:varchar(50);column:id_commission_setting" json:"id"`
	PersentaseKomisi float64   `gorm:"type:decimal(5,2);column:persentase_komisi" json:"rate"`
	BerlakuMulai     time.Time `gorm:"column:berlaku_mulai" json:"startDate"`
	BerlakuSampai    time.Time `gorm:"column:berlaku_sampai" json:"endDate"`
	IsActive         bool      `gorm:"column:is_active" json:"isActive"`
	CreatedBy        string    `gorm:"type:varchar(50);column:created_by" json:"createdBy"`
}