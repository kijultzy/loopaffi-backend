-- Hapus data yang ada jika diperlukan (hati-hati, hapus baris komentar di bawah jika ingin mereset data)
-- TRUNCATE notifications, payments, commissions, sale_items, sales, commission_settings, products, users, roles, payment_methods CASCADE;

-- 1. Insert Roles
INSERT INTO roles (id_role, nama_role) VALUES 
('role_admin', 'Admin'),
('role_affiliate', 'Affiliate')
ON CONFLICT (id_role) DO NOTHING;

-- 2. Insert Users
-- Password untuk semua akun adalah "password123" (sudah di-hash bcrypt)
INSERT INTO users (id_user, id_role, nama_user, email, password_hash, no_hp, status_user, created_at, updated_at) VALUES 
('user_admin1', 'role_admin', 'Admin Utama', 'admin@loopaffi.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '081234567890', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('user_aff1', 'role_affiliate', 'Budi Santoso', 'budi@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '081111111111', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('user_aff2', 'role_affiliate', 'Siti Aminah', 'siti@gmail.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '082222222222', 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id_user) DO NOTHING;

-- 3. Insert Products
INSERT INTO products (id_product, nama_product, sku, harga_default, status_product) VALUES 
('prod_1', 'Kemeja Pria', 'KMJ-001', 150000.00, 'active'),
('prod_2', 'Celana Jeans', 'CLN-001', 250000.00, 'active'),
('prod_3', 'Sepatu Sneakers', 'SPT-001', 350000.00, 'active')
ON CONFLICT (id_product) DO NOTHING;

-- 4. Insert Payment Methods
INSERT INTO payment_methods (id_payment_method, nama_metode) VALUES 
('pay_1', 'Transfer Bank BCA'),
('pay_2', 'Transfer Bank Mandiri'),
('pay_3', 'GoPay'),
('pay_4', 'OVO')
ON CONFLICT (id_payment_method) DO NOTHING;

-- 5. Insert Commission Settings
INSERT INTO commission_settings (id_commission_setting, persentase_komisi, berlaku_mulai, berlaku_sampai, is_active, created_by) VALUES 
('set_1', 10.00, '2023-01-01 00:00:00', '2026-12-31 23:59:59', true, 'user_admin1'),
('set_2', 15.00, '2026-12-01 00:00:00', '2026-12-31 23:59:59', false, 'user_admin1')
ON CONFLICT (id_commission_setting) DO NOTHING;

-- 6. Insert Sales
INSERT INTO sales (id_sale, id_user, tgl_penjualan, total_amount, status_sale) VALUES 
('sale_1', 'user_aff1', '2026-05-10 10:00:00', 400000.00, 'completed'),
('sale_2', 'user_aff1', '2026-05-11 14:30:00', 250000.00, 'completed'),
('sale_3', 'user_aff2', '2026-05-12 09:15:00', 350000.00, 'pending')
ON CONFLICT (id_sale) DO NOTHING;

-- 7. Insert Sale Items
INSERT INTO sale_items (id_sale_item, id_sale, id_product, qty, harga_satuan, subtotal) VALUES 
('item_1', 'sale_1', 'prod_1', 1, 150000.00, 150000.00),
('item_2', 'sale_1', 'prod_2', 1, 250000.00, 250000.00),
('item_3', 'sale_2', 'prod_2', 1, 250000.00, 250000.00),
('item_4', 'sale_3', 'prod_3', 1, 350000.00, 350000.00)
ON CONFLICT (id_sale_item) DO NOTHING;

-- 8. Insert Commissions
INSERT INTO commissions (id_commission, id_sale, id_affiliate, id_commission_setting, jumlah_komisi, tgl_hitung, status_komisi) VALUES 
('comm_1', 'sale_1', 'user_aff1', 'set_1', 40000.00, '2026-05-10 10:05:00', 'pending'),
('comm_2', 'sale_2', 'user_aff1', 'set_1', 25000.00, '2026-05-11 14:35:00', 'paid')
ON CONFLICT (id_commission) DO NOTHING;

-- 9. Insert Payments
INSERT INTO payments (id_payment, id_commission, id_affiliate, id_payment_method, jumlah_bayar, tgl_pembayaran, status_bayar) VALUES 
('paym_1', 'comm_2', 'user_aff1', 'pay_1', 25000.00, '2026-05-11 16:00:00', 'completed')
ON CONFLICT (id_payment) DO NOTHING;

-- 10. Insert Notifications
INSERT INTO notifications (id_notification, id_user, judul, pesan, is_read, created_at) VALUES 
('notif_1', 'user_aff1', 'Penjualan Berhasil', 'Anda mendapatkan komisi baru sebesar Rp 40.000.', false, CURRENT_TIMESTAMP),
('notif_2', 'user_aff1', 'Pembayaran Berhasil', 'Komisi sebesar Rp 25.000 telah ditransfer ke rekening Anda.', true, CURRENT_TIMESTAMP)
ON CONFLICT (id_notification) DO NOTHING;
