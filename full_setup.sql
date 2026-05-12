-- ================================================================
-- LOOPAFFI - FULL DATABASE SETUP
-- ================================================================
-- FILE INI BERISI:
--   BAGIAN A: Membuat tabel (CREATE TABLE)
--   BAGIAN B: Mengisi data dummy (INSERT)
--
-- CARA PAKAI:
--   1. Buka pgAdmin 4
--   2. Klik kanan "Databases" → Create → Database → Nama: loopaffi_db → Save
--   3. Klik kanan database "loopaffi_db" → Query Tool
--   4. Copy-paste SELURUH isi file ini ke Query Tool
--   5. Tekan F5 (Execute)
-- ================================================================


-- =============================================
-- BAGIAN A: BUAT SEMUA TABEL (10 TABEL)
-- =============================================

-- A1. Tabel Roles
CREATE TABLE IF NOT EXISTS roles (
    id_role VARCHAR(50) PRIMARY KEY,
    nama_role VARCHAR(50) NOT NULL
);

-- A2. Tabel Users
CREATE TABLE IF NOT EXISTS users (
    id_user VARCHAR(50) PRIMARY KEY,
    id_role VARCHAR(50) NOT NULL REFERENCES roles(id_role),
    nama_user VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    no_hp VARCHAR(20),
    status_user VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- A3. Tabel Products
CREATE TABLE IF NOT EXISTS products (
    id_product VARCHAR(50) PRIMARY KEY,
    nama_product VARCHAR(100) NOT NULL,
    sku VARCHAR(50),
    harga_default DECIMAL(15,2) NOT NULL,
    status_product VARCHAR(50) DEFAULT 'active'
);

-- A4. Tabel Payment Methods
CREATE TABLE IF NOT EXISTS payment_methods (
    id_payment_method VARCHAR(50) PRIMARY KEY,
    nama_metode VARCHAR(100) NOT NULL
);

-- A5. Tabel Commission Settings
CREATE TABLE IF NOT EXISTS commission_settings (
    id_commission_setting VARCHAR(50) PRIMARY KEY,
    persentase_komisi DECIMAL(5,2) NOT NULL,
    berlaku_mulai TIMESTAMP NOT NULL,
    berlaku_sampai TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_by VARCHAR(50) REFERENCES users(id_user)
);

-- A6. Tabel Sales (Penjualan)
CREATE TABLE IF NOT EXISTS sales (
    id_sale VARCHAR(50) PRIMARY KEY,
    id_user VARCHAR(50) NOT NULL REFERENCES users(id_user),
    tgl_penjualan TIMESTAMP NOT NULL,
    total_amount DECIMAL(15,2) NOT NULL,
    status_sale VARCHAR(50) DEFAULT 'pending'
);

-- A7. Tabel Sale Items (Detail Penjualan)
CREATE TABLE IF NOT EXISTS sale_items (
    id_sale_item VARCHAR(50) PRIMARY KEY,
    id_sale VARCHAR(50) NOT NULL REFERENCES sales(id_sale),
    id_product VARCHAR(50) NOT NULL REFERENCES products(id_product),
    qty INTEGER NOT NULL,
    harga_satuan DECIMAL(15,2) NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL
);

-- A8. Tabel Commissions (Komisi)
CREATE TABLE IF NOT EXISTS commissions (
    id_commission VARCHAR(50) PRIMARY KEY,
    id_sale VARCHAR(50) NOT NULL REFERENCES sales(id_sale),
    id_affiliate VARCHAR(50) NOT NULL REFERENCES users(id_user),
    id_commission_setting VARCHAR(50) REFERENCES commission_settings(id_commission_setting),
    jumlah_komisi DECIMAL(15,2) NOT NULL,
    tgl_hitung TIMESTAMP NOT NULL,
    status_komisi VARCHAR(50) DEFAULT 'pending'
);

-- A9. Tabel Payments (Pembayaran)
CREATE TABLE IF NOT EXISTS payments (
    id_payment VARCHAR(50) PRIMARY KEY,
    id_commission VARCHAR(50) NOT NULL REFERENCES commissions(id_commission),
    id_affiliate VARCHAR(50) NOT NULL REFERENCES users(id_user),
    id_payment_method VARCHAR(50) REFERENCES payment_methods(id_payment_method),
    jumlah_bayar DECIMAL(15,2) NOT NULL,
    tgl_pembayaran TIMESTAMP NOT NULL,
    status_bayar VARCHAR(50) DEFAULT 'pending'
);

-- A10. Tabel Notifications (Notifikasi)
CREATE TABLE IF NOT EXISTS notifications (
    id_notification VARCHAR(50) PRIMARY KEY,
    id_user VARCHAR(50) NOT NULL REFERENCES users(id_user),
    judul VARCHAR(100),
    pesan TEXT,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- =============================================
-- BAGIAN B: ISI DATA DUMMY
-- =============================================

-- B1. Roles
INSERT INTO roles (id_role, nama_role) VALUES
('ROLE-001', 'Admin'),
('ROLE-002', 'Affiliate')
ON CONFLICT (id_role) DO NOTHING;

-- B2. Users (password = "password123", bcrypt hash)
INSERT INTO users (id_user, id_role, nama_user, email, password_hash, no_hp, status_user, created_at, updated_at) VALUES
('1', 'ROLE-001', 'Admin LoopAffi', 'admin@loopaffi.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 '081234567890', 'active', NOW(), NOW()),

('2', 'ROLE-002', 'Budi Santoso', 'budi@gmail.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 '081111111111', 'active', NOW(), NOW()),

('3', 'ROLE-002', 'Siti Aminah', 'siti@gmail.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 '082222222222', 'active', NOW(), NOW()),

('4', 'ROLE-002', 'Rizky Pratama', 'rizky@gmail.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 '083333333333', 'active', NOW(), NOW()),

('5', 'ROLE-002', 'Diana Putri', 'diana@gmail.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
 '084444444444', 'active', NOW(), NOW())
ON CONFLICT (id_user) DO NOTHING;

-- B3. Products
INSERT INTO products (id_product, nama_product, sku, harga_default, status_product) VALUES
('PROD-001', 'LoopAffi Basic Plan',   'LA-BASIC-01', 500000.00,  'active'),
('PROD-002', 'LoopAffi Pro Plan',     'LA-PRO-01',   1500000.00, 'active'),
('PROD-003', 'LoopAffi Enterprise',   'LA-ENT-01',   5000000.00, 'active'),
('PROD-004', 'LoopAffi Starter Pack', 'LA-START-01', 250000.00,  'active')
ON CONFLICT (id_product) DO NOTHING;

-- B4. Payment Methods
INSERT INTO payment_methods (id_payment_method, nama_metode) VALUES
('PM-001', 'Transfer Bank BCA'),
('PM-002', 'Transfer Bank Mandiri'),
('PM-003', 'GoPay'),
('PM-004', 'OVO')
ON CONFLICT (id_payment_method) DO NOTHING;

-- B5. Commission Settings (10% aktif)
INSERT INTO commission_settings (id_commission_setting, persentase_komisi, berlaku_mulai, berlaku_sampai, is_active, created_by) VALUES
('CS-001', 10.00, '2026-01-01', '2027-12-31', true,  '1'),
('CS-002', 15.00, '2027-01-01', '2027-12-31', false, '1')
ON CONFLICT (id_commission_setting) DO NOTHING;

-- B6. Sales (Penjualan)
INSERT INTO sales (id_sale, id_user, tgl_penjualan, total_amount, status_sale) VALUES
('SALE-001', '2', '2026-05-01 10:30:00', 500000.00,  'completed'),
('SALE-002', '2', '2026-05-03 14:00:00', 1500000.00, 'completed'),
('SALE-003', '3', '2026-05-05 09:15:00', 5000000.00, 'completed'),
('SALE-004', '3', '2026-05-07 16:45:00', 500000.00,  'completed'),
('SALE-005', '4', '2026-05-08 11:00:00', 250000.00,  'completed'),
('SALE-006', '4', '2026-05-09 13:30:00', 1500000.00, 'pending'),
('SALE-007', '5', '2026-05-10 08:00:00', 5000000.00, 'pending')
ON CONFLICT (id_sale) DO NOTHING;

-- B7. Sale Items (Detail per Penjualan)
INSERT INTO sale_items (id_sale_item, id_sale, id_product, qty, harga_satuan, subtotal) VALUES
('SI-001', 'SALE-001', 'PROD-001', 1, 500000.00,  500000.00),
('SI-002', 'SALE-002', 'PROD-002', 1, 1500000.00, 1500000.00),
('SI-003', 'SALE-003', 'PROD-003', 1, 5000000.00, 5000000.00),
('SI-004', 'SALE-004', 'PROD-001', 1, 500000.00,  500000.00),
('SI-005', 'SALE-005', 'PROD-004', 1, 250000.00,  250000.00),
('SI-006', 'SALE-006', 'PROD-002', 1, 1500000.00, 1500000.00),
('SI-007', 'SALE-007', 'PROD-003', 1, 5000000.00, 5000000.00)
ON CONFLICT (id_sale_item) DO NOTHING;

-- B8. Commissions (Komisi — 10% dari total penjualan completed)
INSERT INTO commissions (id_commission, id_sale, id_affiliate, id_commission_setting, jumlah_komisi, tgl_hitung, status_komisi) VALUES
('COMM-001', 'SALE-001', '2', 'CS-001', 50000.00,  '2026-05-01 10:35:00', 'paid'),
('COMM-002', 'SALE-002', '2', 'CS-001', 150000.00, '2026-05-03 14:05:00', 'paid'),
('COMM-003', 'SALE-003', '3', 'CS-001', 500000.00, '2026-05-05 09:20:00', 'pending'),
('COMM-004', 'SALE-004', '3', 'CS-001', 50000.00,  '2026-05-07 16:50:00', 'pending'),
('COMM-005', 'SALE-005', '4', 'CS-001', 25000.00,  '2026-05-08 11:05:00', 'pending')
ON CONFLICT (id_commission) DO NOTHING;

-- B9. Payments (Pembayaran komisi yang sudah dibayar)
INSERT INTO payments (id_payment, id_commission, id_affiliate, id_payment_method, jumlah_bayar, tgl_pembayaran, status_bayar) VALUES
('PAY-001', 'COMM-001', '2', 'PM-001', 50000.00,  '2026-05-02 09:00:00', 'completed'),
('PAY-002', 'COMM-002', '2', 'PM-003', 150000.00, '2026-05-04 10:00:00', 'completed')
ON CONFLICT (id_payment) DO NOTHING;

-- B10. Notifications
INSERT INTO notifications (id_notification, id_user, judul, pesan, is_read, created_at) VALUES
('NOTIF-001', '2', 'Komisi Dibayar', 'Komisi Rp 50.000 dari penjualan SALE-001 telah ditransfer ke BCA Anda.', true, '2026-05-02 09:01:00'),
('NOTIF-002', '2', 'Komisi Dibayar', 'Komisi Rp 150.000 dari penjualan SALE-002 telah ditransfer via GoPay.', true, '2026-05-04 10:01:00'),
('NOTIF-003', '3', 'Komisi Baru', 'Anda mendapatkan komisi Rp 500.000 dari penjualan SALE-003. Menunggu pembayaran.', false, '2026-05-05 09:21:00'),
('NOTIF-004', '3', 'Komisi Baru', 'Anda mendapatkan komisi Rp 50.000 dari penjualan SALE-004. Menunggu pembayaran.', false, '2026-05-07 16:51:00'),
('NOTIF-005', '4', 'Penjualan Tercatat', 'Penjualan SALE-005 senilai Rp 250.000 berhasil dicatat.', false, '2026-05-08 11:06:00')
ON CONFLICT (id_notification) DO NOTHING;


-- ================================================================
-- SELESAI! Verifikasi dengan query di bawah:
-- ================================================================
SELECT 'roles' AS tabel, COUNT(*) AS jumlah FROM roles
UNION ALL SELECT 'users', COUNT(*) FROM users
UNION ALL SELECT 'products', COUNT(*) FROM products
UNION ALL SELECT 'payment_methods', COUNT(*) FROM payment_methods
UNION ALL SELECT 'commission_settings', COUNT(*) FROM commission_settings
UNION ALL SELECT 'sales', COUNT(*) FROM sales
UNION ALL SELECT 'sale_items', COUNT(*) FROM sale_items
UNION ALL SELECT 'commissions', COUNT(*) FROM commissions
UNION ALL SELECT 'payments', COUNT(*) FROM payments
UNION ALL SELECT 'notifications', COUNT(*) FROM notifications
ORDER BY tabel;
