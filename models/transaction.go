package models

import "time"

// 1. Model untuk tabel sales (Penjualan)
type Sale struct {
	ID           string    `gorm:"primaryKey;type:varchar(50);column:id_sale" json:"id"`
	UserID       string    `gorm:"type:varchar(50);column:id_user" json:"affiliateId"`
	TglPenjualan time.Time `gorm:"column:tgl_penjualan" json:"date"`
	TotalAmount  float64   `gorm:"type:decimal(15,2);column:total_amount" json:"amount"`
	StatusSale   string    `gorm:"type:varchar(50);column:status_sale" json:"status"`

	// Relasi
	User      User       `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	SaleItems []SaleItem `gorm:"foreignKey:SaleID;references:ID" json:"items,omitempty"`
}

// 2. Model untuk tabel sale_items (Detail Penjualan)
type SaleItem struct {
	ID          string  `gorm:"primaryKey;type:varchar(50);column:id_sale_item" json:"id"`
	SaleID      string  `gorm:"type:varchar(50);column:id_sale" json:"saleId"`
	ProductID   string  `gorm:"type:varchar(50);column:id_product" json:"productId"`
	Qty         int     `gorm:"column:qty" json:"qty"`
	HargaSatuan float64 `gorm:"type:decimal(15,2);column:harga_satuan" json:"unitPrice"`
	Subtotal    float64 `gorm:"type:decimal(15,2);column:subtotal" json:"subtotal"`

	// Relasi
	Sale    Sale    `gorm:"foreignKey:SaleID;references:ID" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

// 3. Model untuk tabel commissions (Komisi)
type Commission struct {
	ID                  string    `gorm:"primaryKey;type:varchar(50);column:id_commission" json:"id"`
	SaleID              string    `gorm:"type:varchar(50);column:id_sale" json:"saleId"`
	AffiliateID         string    `gorm:"type:varchar(50);column:id_affiliate" json:"affiliateId"`
	CommissionSettingID string    `gorm:"type:varchar(50);column:id_commission_setting" json:"commissionSettingId"`
	JumlahKomisi        float64   `gorm:"type:decimal(15,2);column:jumlah_komisi" json:"amount"`
	TglHitung           time.Time `gorm:"column:tgl_hitung" json:"date"`
	StatusKomisi        string    `gorm:"type:varchar(50);column:status_komisi" json:"status"`

	// Relasi
	Sale              Sale              `gorm:"foreignKey:SaleID;references:ID" json:"sale,omitempty"`
	CommissionSetting CommissionSetting `gorm:"foreignKey:CommissionSettingID;references:ID" json:"commissionSetting,omitempty"`
	User              User              `gorm:"foreignKey:AffiliateID;references:ID" json:"affiliate,omitempty"`
}

// 4. Model untuk tabel payments (Pembayaran)
type Payment struct {
	ID              string    `gorm:"primaryKey;type:varchar(50);column:id_payment" json:"id"`
	CommissionID    string    `gorm:"type:varchar(50);column:id_commission" json:"commissionId"`
	AffiliateID     string    `gorm:"type:varchar(50);column:id_affiliate" json:"affiliateId"`
	PaymentMethodID string    `gorm:"type:varchar(50);column:id_payment_method" json:"paymentMethodId"`
	JumlahBayar     float64   `gorm:"type:decimal(15,2);column:jumlah_bayar" json:"amount"`
	TglPembayaran   time.Time `gorm:"column:tgl_pembayaran" json:"date"`
	StatusBayar     string    `gorm:"type:varchar(50);column:status_bayar" json:"status"`

	// Relasi
	Commission    Commission    `gorm:"foreignKey:CommissionID;references:ID" json:"commission,omitempty"`
	PaymentMethod PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID" json:"paymentMethod,omitempty"`
	User          User          `gorm:"foreignKey:AffiliateID;references:ID" json:"affiliate,omitempty"`
}

// 5. Model untuk tabel notifications (Notifikasi)
type Notification struct {
	ID        string    `gorm:"primaryKey;type:varchar(50);column:id_notification" json:"id"`
	UserID    string    `gorm:"type:varchar(50);column:id_user" json:"userId"`
	Judul     string    `gorm:"type:varchar(100);column:judul" json:"title"`
	Pesan     string    `gorm:"type:text;column:pesan" json:"message"`
	IsRead    bool      `gorm:"default:false;column:is_read" json:"read"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"date"`

	// Relasi
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}