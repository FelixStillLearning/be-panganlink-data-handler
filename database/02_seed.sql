-- Insert Master Data for Commodities
INSERT INTO komoditas (id, nama, satuan, kategori) VALUES
('kmdB-001', 'Beras Pandan Wangi', 'kg', 'Beras'),
('kmdC-001', 'Cabai Merah Keriting', 'kg', 'Sayuran'),
('kmdBM-001', 'Bawang Merah', 'kg', 'Sayuran'),
('kmdT-001', 'Tomat', 'kg', 'Sayuran');

-- Insert default Admin user
INSERT INTO users (name, email, password, role, location) VALUES
('Super Admin', 'admin@panganlink.com', '$2a$12$N9/mX0pG2U99z1Q3L8M7..Ue0/2E2T9W9C5v3Wf.9QZ/2Q.4y8O5C', 'admin', 'Jakarta');

-- Insert initial dummy data for harga_pasar
INSERT INTO harga_pasar (komoditas_id, harga, wilayah, tanggal) VALUES
('kmdB-001', 14000, 'Nasional', '2026-06-01'),
('kmdB-001', 14200, 'Nasional', '2026-06-02'),
('kmdB-001', 14100, 'Nasional', '2026-06-03'),
('kmdB-001', 14300, 'Nasional', '2026-06-04'),
('kmdB-001', 14300, 'Nasional', '2026-06-05');
