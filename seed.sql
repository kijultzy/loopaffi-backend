-- ============================================
-- LOOPAFFI - SQL SEED DATA MASTER
-- Jalankan setelah go run main.go (auto-migrate)
-- ============================================

-- 1. SEED ROLES
INSERT INTO roles (id_role, nama_role) VALUES
('ROLE-001', 'Admin'),
('ROLE-002', 'Affiliate');

-- 2. SEED USERS (password = "password")
-- ID sesuai mockUsers di frontend store.ts
INSERT INTO users (id_user, id_role, nama_user, email, password_hash, no_hp, status_user, created_at, updated_at) VALUES
('1', 'ROLE-001', 'Admin User', 'admin@loopaffi.com', 'password', '081234567890', 'active', NOW(), NOW()),
('2', 'ROLE-002', 'John Affiliate', 'john@example.com', 'password', '081234567891', 'active', NOW(), NOW()),
('3', 'ROLE-002', 'Jane Marketer', 'jane@example.com', 'password', '081234567892', 'active', NOW(), NOW());

-- 3. SEED PRODUCTS
INSERT INTO products (id_product, nama_product, sku, harga_default, status_product) VALUES
('PROD-001', 'LoopAffi Basic Plan', 'LA-BASIC-01', 500000.00, 'active'),
('PROD-002', 'LoopAffi Pro Plan', 'LA-PRO-01', 1500000.00, 'active'),
('PROD-003', 'LoopAffi Enterprise', 'LA-ENT-01', 5000000.00, 'active');

-- 4. SEED PAYMENT METHODS
INSERT INTO payment_methods (id_payment_method, nama_metode) VALUES
('PM-001', 'Transfer Bank BCA'),
('PM-002', 'Transfer Bank Mandiri'),
('PM-003', 'GoPay'),
('PM-004', 'OVO');

-- 5. SEED COMMISSION SETTINGS (10% aktif)
-- Sesuai globalCommissionRate: 0.1 di frontend
INSERT INTO commission_settings (id_commission_setting, persentase_komisi, berlaku_mulai, berlaku_sampai, is_active, created_by) VALUES
('CS-001', 10.00, '2026-01-01', '2027-12-31', true, '1');
