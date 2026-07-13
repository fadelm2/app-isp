-- Clean up existing dummy seed data if any (we can keep original noc1/admin1, but clear others to start clean)
DELETE FROM payments WHERE invoice_id LIKE 'INV-SEED-%';
DELETE FROM invoices WHERE id LIKE 'INV-SEED-%';
DELETE FROM customer_histories WHERE customer_id LIKE 'CUST-SEED-%';
DELETE FROM customers WHERE id LIKE 'CUST-SEED-%';
DELETE FROM registrations WHERE id LIKE 'REG-SEED-%';
DELETE FROM users WHERE id LIKE 'user-seed-%';

-- Ensure we have packages
INSERT IGNORE INTO internet_packages (id, name, speed_mbps, price, installation_fee, tax_rate, is_active, created_at, updated_at) VALUES 
('pkg-seed-10m', 'Greenet Lite 10M', 10, 125000, 150000, 0.11, 1, 1783899129160, 1783899129160),
('pkg-seed-20m', 'Greenet Family 20M', 20, 185000, 150000, 0.11, 1, 1783899129160, 1783899129160),
('pkg-seed-50m', 'Greenet Pro 50M', 50, 325000, 250000, 0.11, 1, 1783899129160, 1783899129160),
('pkg-seed-100m', 'Greenet Gamer 100M', 100, 550000, 250000, 0.11, 1, 1783899129160, 1783899129160);

-- Insert 10 Registrations (Some pending, some surveying, some rejected)
INSERT INTO registrations (id, full_name, nik, birth_place, birth_date, gender, email, phone, installation_address, billing_address, package_id, latitude, longitude, status, ktp_path, province, city, district, village, created_at, updated_at) VALUES
('REG-SEED-1', 'Budi Santoso', '3201234567890001', 'Jakarta', '1990-05-12', 'Laki-laki', 'budi.santoso@gmail.com', '081234567890', 'Jl. Sudirman No. 12', 'Jl. Sudirman No. 12', 'pkg-seed-20m', -6.2088, 106.8456, 'pending', '/storage/uploads/ktp/dummy.jpg', 'DKI Jakarta', 'Jakarta Pusat', 'Tanah Abang', 'Kebon Melati', 1783899129160, 1783899129160),
('REG-SEED-2', 'Siti Aminah', '3201234567890002', 'Bandung', '1993-08-22', 'Perempuan', 'siti.aminah@gmail.com', '081234567891', 'Jl. Dago No. 45', 'Jl. Dago No. 45', 'pkg-seed-10m', -6.9175, 107.6191, 'surveying', '/storage/uploads/ktp/dummy.jpg', 'Jawa Barat', 'Kota Bandung', 'Coblong', 'Dago', 1783899229160, 1783899229160),
('REG-SEED-3', 'Dewi Lestari', '3201234567890003', 'Surabaya', '1988-11-05', 'Perempuan', 'dewi.lestari@gmail.com', '081234567892', 'Jl. Basuki Rahmat No. 78', 'Jl. Basuki Rahmat No. 78', 'pkg-seed-50m', -7.2575, 112.7521, 'rejected', '/storage/uploads/ktp/dummy.jpg', 'Jawa Timur', 'Kota Surabaya', 'Genteng', 'Embong Kaliasin', 1783899329160, 1783899329160),
('REG-SEED-4', 'Hendra Wijaya', '3201234567890004', 'Semarang', '1995-02-14', 'Laki-laki', 'hendra.wijaya@gmail.com', '081234567893', 'Jl. Pemuda No. 101', 'Jl. Pemuda No. 101', 'pkg-seed-20m', -6.9932, 110.4203, 'pending', '/storage/uploads/ktp/dummy.jpg', 'Jawa Tengah', 'Kota Semarang', 'Semarang Tengah', 'Sekayu', 1783899429160, 1783899429160),
('REG-SEED-5', 'Rian Hidayat', '3201234567890005', 'Yogyakarta', '1991-07-30', 'Laki-laki', 'rian.hidayat@gmail.com', '081234567894', 'Jl. Malioboro No. 56', 'Jl. Malioboro No. 56', 'pkg-seed-10m', -7.7956, 110.3695, 'surveying', '/storage/uploads/ktp/dummy.jpg', 'DI Yogyakarta', 'Kota Yogyakarta', 'Gedongtengen', 'Sosromenduran', 1783899529160, 1783899529160),
('REG-SEED-6', 'Anisa Putri', '3201234567890006', 'Malang', '1994-04-18', 'Perempuan', 'anisa.putri@gmail.com', '081234567895', 'Jl. Ijen No. 12', 'Jl. Ijen No. 12', 'pkg-seed-50m', -7.9839, 112.6214, 'pending', '/storage/uploads/ktp/dummy.jpg', 'Jawa Timur', 'Kota Malang', 'Klojen', 'Oro-oro Dowo', 1783899629160, 1783899629160),
('REG-SEED-7', 'Rudi Hermawan', '3201234567890007', 'Medan', '1987-12-25', 'Laki-laki', 'rudi.hermawan@gmail.com', '081234567896', 'Jl. Gajah Mada No. 89', 'Jl. Gajah Mada No. 89', 'pkg-seed-100m', 3.5952, 98.6722, 'rejected', '/storage/uploads/ktp/dummy.jpg', 'Sumatera Utara', 'Kota Medan', 'Medan Petisah', 'Petisah Tengah', 1783899729160, 1783899729160),
('REG-SEED-8', 'Eka Saputra', '3201234567890008', 'Palembang', '1992-09-09', 'Laki-laki', 'eka.saputra@gmail.com', '081234567897', 'Jl. Sudirman No. 34', 'Jl. Sudirman No. 34', 'pkg-seed-20m', -2.9909, 104.7567, 'pending', '/storage/uploads/ktp/dummy.jpg', 'Sumatera Selatan', 'Kota Palembang', 'Ilir Barat I', 'Bukit Baru', 1783899829160, 1783899829160),
('REG-SEED-9', 'Mega Utami', '3201234567890009', 'Makassar', '1996-06-15', 'Perempuan', 'mega.utami@gmail.com', '081234567898', 'Jl. Pettarani No. 112', 'Jl. Pettarani No. 112', 'pkg-seed-10M', -5.1477, 119.4327, 'surveying', '/storage/uploads/ktp/dummy.jpg', 'Sulawesi Selatan', 'Kota Makassar', 'Rappocini', 'Gunung Sari', 1783899929160, 1783899929160),
('REG-SEED-10', 'Aditya Pratama', '3201234567890010', 'Denpasar', '1989-10-10', 'Laki-laki', 'aditya.pratama@gmail.com', '081234567899', 'Jl. Teuku Umar No. 202', 'Jl. Teuku Umar No. 202', 'pkg-seed-50m', -8.6705, 115.2126, 'pending', '/storage/uploads/ktp/dummy.jpg', 'Bali', 'Kota Denpasar', 'Denpasar Barat', 'Pemecutan', 1783900029160, 1783900029160);

-- Insert 10 Approved Registrations that became Active/Suspended Customers
-- First insert users
INSERT INTO users (id, role_id, username, email, password, company_name, created_at, updated_at) VALUES
('user-seed-11', 2, 'Geri Maulana', 'geri.maulana@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900129160, 1783900129160),
('user-seed-12', 2, 'Laras Ati', 'laras.ati@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900229160, 1783900229160),
('user-seed-13', 2, 'Fajar Ramadhan', 'fajar.ramadhan@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900329160, 1783900329160),
('user-seed-14', 2, 'Novianti', 'novianti@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900429160, 1783900429160),
('user-seed-15', 2, 'Reza Rahadian', 'reza.rahadian@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900529160, 1783900529160),
('user-seed-16', 2, 'Fitri Handayani', 'fitri.handayani@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900629160, 1783900629160),
('user-seed-17', 2, 'Dimas Anggara', 'dimas.anggara@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900729160, 1783900729160),
('user-seed-18', 2, 'Yulia Citra', 'yulia.citra@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900829160, 1783900829160),
('user-seed-19', 2, 'Kevin Sanjaya', 'kevin.sanjaya@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783900929160, 1783900929160),
('user-seed-20', 2, 'Gita Gutawa', 'gita.gutawa@gmail.com', '$2a$10$Z4MR5mDWzrDxVCCasdu5VeTf5DbYcsyMb/aMeP4BlDFoOeLO2.R9y', 'GREENET', 1783901029160, 1783901029160);

-- Insert registrations for these approved ones
INSERT INTO registrations (id, full_name, nik, birth_place, birth_date, gender, email, phone, installation_address, billing_address, package_id, latitude, longitude, status, ktp_path, province, city, district, village, created_at, updated_at) VALUES
('REG-SEED-11', 'Geri Maulana', '3201234567890011', 'Surakarta', '1993-11-20', 'Laki-laki', 'geri.maulana@gmail.com', '081234567911', 'Jl. Slamet Riyadi No. 45', 'Jl. Slamet Riyadi No. 45', 'pkg-seed-20m', -7.5684, 110.8219, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Jawa Tengah', 'Kota Surakarta', 'Laweyan', 'Pajang', 1783900129160, 1783900129160),
('REG-SEED-12', 'Laras Ati', '3201234567890012', 'Balikpapan', '1991-03-12', 'Perempuan', 'laras.ati@gmail.com', '081234567912', 'Jl. Sudirman No. 78', 'Jl. Sudirman No. 78', 'pkg-seed-10m', -1.2654, 116.8312, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Kalimantan Timur', 'Kota Balikpapan', 'Balikpapan Kota', 'Prapatan', 1783900229160, 1783900229160),
('REG-SEED-13', 'Fajar Ramadhan', '3201234567890013', 'Banjarmasin', '1994-07-28', 'Laki-laki', 'fajar.ramadhan@gmail.com', '081234567913', 'Jl. Ahmad Yani No. 12', 'Jl. Ahmad Yani No. 12', 'pkg-seed-50m', -3.3186, 114.5944, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Kalimantan Selatan', 'Kota Banjarmasin', 'Banjarmasin Tengah', 'Telaga Biru', 1783900329160, 1783900329160),
('REG-SEED-14', 'Novianti', '3201234567890014', 'Manado', '1988-12-15', 'Perempuan', 'novianti@gmail.com', '081234567914', 'Jl. Sam Ratulangi No. 56', 'Jl. Sam Ratulangi No. 56', 'pkg-seed-20m', 1.4748, 124.8428, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Sulawesi Utara', 'Kota Manado', 'Wenang', 'Wenang Utara', 1783900429160, 1783900429160),
('REG-SEED-15', 'Reza Rahadian', '3201234567890015', 'Bogor', '1990-10-09', 'Laki-laki', 'reza.rahadian@gmail.com', '081234567915', 'Jl. Pajajaran No. 99', 'Jl. Pajajaran No. 99', 'pkg-seed-100m', -6.5971, 106.7986, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Jawa Barat', 'Kota Bogor', 'Bogor Tengah', 'Babakan', 1783900529160, 1783900529160),
('REG-SEED-16', 'Fitri Handayani', '3201234567890016', 'Cirebon', '1995-04-03', 'Perempuan', 'fitri.handayani@gmail.com', '081234567916', 'Jl. Kartini No. 34', 'Jl. Kartini No. 34', 'pkg-seed-10m', -6.7216, 108.5562, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Jawa Barat', 'Kota Cirebon', 'Kejaksan', 'Kesenden', 1783900629160, 1783900629160),
('REG-SEED-17', 'Dimas Anggara', '3201234567890017', 'Pekanbaru', '1992-05-24', 'Laki-laki', 'dimas.anggara@gmail.com', '081234567917', 'Jl. Sudirman No. 120', 'Jl. Sudirman No. 120', 'pkg-seed-20m', 0.5071, 101.4478, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Riau', 'Kota Pekanbaru', 'Pekanbaru Kota', 'Simpang Empat', 1783900729160, 1783900729160),
('REG-SEED-18', 'Yulia Citra', '3201234567890018', 'Batam', '1996-01-30', 'Perempuan', 'yulia.citra@gmail.com', '081234567918', 'Jl. Nagoya No. 4', 'Jl. Nagoya No. 4', 'pkg-seed-50m', 1.1443, 104.0076, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Kepulauan Riau', 'Kota Batam', 'Lubuk Baja', 'Lubuk Baja Kota', 1783900829160, 1783900829160),
('REG-SEED-19', 'Kevin Sanjaya', '3201234567890019', 'Cianjur', '1993-02-18', 'Laki-laki', 'kevin.sanjaya@gmail.com', '081234567919', 'Jl. Raya Cianjur No. 88', 'Jl. Raya Cianjur No. 88', 'pkg-seed-100m', -6.8208, 107.1417, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Jawa Barat', 'Kabupaten Cianjur', 'Cianjur', 'Pamoyanan', 1783900929160, 1783900929160),
('REG-SEED-20', 'Gita Gutawa', '3201234567890020', 'Serang', '1994-09-04', 'Perempuan', 'gita.gutawa@gmail.com', '081234567920', 'Jl. Veteran No. 15', 'Jl. Veteran No. 15', 'pkg-seed-20m', -6.1153, 106.1503, 'approved', '/storage/uploads/ktp/dummy.jpg', 'Banten', 'Kota Serang', 'Serang', 'Cimuncang', 1783901029160, 1783901029160);

-- Insert customers
INSERT INTO customers (id, registration_id, user_id, status, package_id, router_id, ppp_username, ppp_password, radius_username, radius_password, due_date_day, created_at, updated_at) VALUES
('CUST-SEED-11', 'REG-SEED-11', 'user-seed-11', 'active', 'pkg-seed-20m', NULL, 'Geri Maulana@greenet', 'pass11', 'Geri Maulana@greenet', 'pass11', 15, 1783900129160, 1783900129160),
('CUST-SEED-12', 'REG-SEED-12', 'user-seed-12', 'active', 'pkg-seed-10m', NULL, 'Laras Ati@greenet', 'pass12', 'Laras Ati@greenet', 'pass12', 15, 1783900229160, 1783900229160),
('CUST-SEED-13', 'REG-SEED-13', 'user-seed-13', 'suspended', 'pkg-seed-50m', NULL, 'Fajar Ramadhan@greenet', 'pass13', 'Fajar Ramadhan@greenet', 'pass13', 15, 1783900329160, 1783900329160),
('CUST-SEED-14', 'REG-SEED-14', 'user-seed-14', 'active', 'pkg-seed-20m', NULL, 'Novianti@greenet', 'pass14', 'Novianti@greenet', 'pass14', 15, 1783900429160, 1783900429160),
('CUST-SEED-15', 'REG-SEED-15', 'user-seed-15', 'active', 'pkg-seed-100m', NULL, 'Reza Rahadian@greenet', 'pass15', 'Reza Rahadian@greenet', 'pass15', 15, 1783900529160, 1783900529160),
('CUST-SEED-16', 'REG-SEED-16', 'user-seed-16', 'active', 'pkg-seed-10m', NULL, 'Fitri Handayani@greenet', 'pass16', 'Fitri Handayani@greenet', 'pass16', 15, 1783900629160, 1783900629160),
('CUST-SEED-17', 'REG-SEED-17', 'user-seed-17', 'suspended', 'pkg-seed-20m', NULL, 'Dimas Anggara@greenet', 'pass17', 'Dimas Anggara@greenet', 'pass17', 15, 1783900729160, 1783900729160),
('CUST-SEED-18', 'REG-SEED-18', 'user-seed-18', 'active', 'pkg-seed-50m', NULL, 'Yulia Citra@greenet', 'pass18', 'Yulia Citra@greenet', 'pass18', 15, 1783900829160, 1783900829160),
('CUST-SEED-19', 'REG-SEED-19', 'user-seed-19', 'active', 'pkg-seed-100m', NULL, 'Kevin Sanjaya@greenet', 'pass19', 'Kevin Sanjaya@greenet', 'pass19', 15, 1783900929160, 1783900929160),
('CUST-SEED-20', 'REG-SEED-20', 'user-seed-20', 'active', 'pkg-seed-20m', NULL, 'Gita Gutawa@greenet', 'pass20', 'Gita Gutawa@greenet', 'pass20', 15, 1783901029160, 1783901029160);

-- Insert Invoices
INSERT INTO invoices (id, customer_id, due_date, period_month, period_year, amount, tax_amount, installation_fee, total_amount, status, snap_token, paid_at, created_at, updated_at) VALUES
('INV-SEED-11', 'CUST-SEED-11', 1783902929160, 7, 2026, 185000, 20350, 150000, 355350, 'paid', 'snap-token-11', 1783900229160, 1783900129160, 1783900229160),
('INV-SEED-12', 'CUST-SEED-12', 1783903029160, 7, 2026, 125000, 13750, 150000, 288750, 'paid', 'snap-token-12', 1783900329160, 1783900229160, 1783900329160),
('INV-SEED-13', 'CUST-SEED-13', 1783903129160, 7, 2026, 325000, 35750, 250000, 610750, 'pending', 'snap-token-13', NULL, 1783900329160, 1783900329160),
('INV-SEED-14', 'CUST-SEED-14', 1783903229160, 7, 2026, 185000, 20350, 150000, 355350, 'paid', 'snap-token-14', 1783900529160, 1783900429160, 1783900529160),
('INV-SEED-15', 'CUST-SEED-15', 1783903329160, 7, 2026, 550000, 60500, 250000, 860500, 'paid', 'snap-token-15', 1783900629160, 1783900529160, 1783900629160),
('INV-SEED-16', 'CUST-SEED-16', 1783903429160, 7, 2026, 125000, 13750, 150000, 288750, 'paid', 'snap-token-16', 1783900729160, 1783900629160, 1783900729160),
('INV-SEED-17', 'CUST-SEED-17', 1783903529160, 7, 2026, 185000, 20350, 150000, 355350, 'pending', 'snap-token-17', NULL, 1783900729160, 1783900729160),
('INV-SEED-18', 'CUST-SEED-18', 1783903629160, 7, 2026, 325000, 35750, 250000, 610750, 'paid', 'snap-token-18', 1783900929160, 1783900829160, 1783900929160),
('INV-SEED-19', 'CUST-SEED-19', 1783903729160, 7, 2026, 550000, 60500, 250000, 860500, 'paid', 'snap-token-19', 1783901029160, 1783900929160, 1783901029160),
('INV-SEED-20', 'CUST-SEED-20', 1783903829160, 7, 2026, 185000, 20350, 150000, 355350, 'paid', 'snap-token-20', 1783901129160, 1783901029160, 1783901129160);

-- Insert Payments
INSERT INTO payments (id, invoice_id, transaction_id, payment_type, paid_amount, status, paid_at, raw_response) VALUES
('PAY-SEED-11', 'INV-SEED-11', 'midtrans-tx-11', 'bank_transfer', 355350, 'settlement', 1783900229160, '{}'),
('PAY-SEED-12', 'INV-SEED-12', 'midtrans-tx-12', 'credit_card', 288750, 'settlement', 1783900329160, '{}'),
('PAY-SEED-14', 'INV-SEED-14', 'midtrans-tx-14', 'gopay', 355350, 'settlement', 1783900529160, '{}'),
('PAY-SEED-15', 'INV-SEED-15', 'midtrans-tx-15', 'bank_transfer', 860500, 'settlement', 1783900629160, '{}'),
('PAY-SEED-16', 'INV-SEED-16', 'midtrans-tx-16', 'shopeepay', 288750, 'settlement', 1783900729160, '{}'),
('PAY-SEED-18', 'INV-SEED-18', 'midtrans-tx-18', 'credit_card', 610750, 'settlement', 1783900929160, '{}'),
('PAY-SEED-19', 'INV-SEED-19', 'midtrans-tx-19', 'bank_transfer', 860500, 'settlement', 1783901029160, '{}'),
('PAY-SEED-20', 'INV-SEED-20', 'midtrans-tx-20', 'gopay', 355350, 'settlement', 1783901129160, '{}');

-- Insert Customer Histories
INSERT INTO customer_histories (id, customer_id, action, notes, created_by, created_at) VALUES
('HIST-SEED-11', 'CUST-SEED-11', 'created', 'Customer registered and approved', 'admin1', 1783900129160),
('HIST-SEED-12', 'CUST-SEED-12', 'created', 'Customer registered and approved', 'admin1', 1783900229160),
('HIST-SEED-13', 'CUST-SEED-13', 'created', 'Customer registered and approved', 'admin1', 1783900329160),
('HIST-SEED-14', 'CUST-SEED-14', 'created', 'Customer registered and approved', 'admin1', 1783900429160),
('HIST-SEED-15', 'CUST-SEED-15', 'created', 'Customer registered and approved', 'admin1', 1783900529160),
('HIST-SEED-16', 'CUST-SEED-16', 'created', 'Customer registered and approved', 'admin1', 1783900629160),
('HIST-SEED-17', 'CUST-SEED-17', 'created', 'Customer registered and approved', 'admin1', 1783900729160),
('HIST-SEED-18', 'CUST-SEED-18', 'created', 'Customer registered and approved', 'admin1', 1783900829160),
('HIST-SEED-19', 'CUST-SEED-19', 'created', 'Customer registered and approved', 'admin1', 1783900929160),
('HIST-SEED-20', 'CUST-SEED-20', 'created', 'Customer registered and approved', 'admin1', 1783901029160);
