-- users: semua role dalam satu tabel, dibedakan via kolom role
CREATE TABLE users (
    id          VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,           -- bcrypt hash
    role        VARCHAR(20) NOT NULL CHECK (role IN ('petani', 'pembeli', 'admin')),
    location    VARCHAR(100),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- komoditas: master data jenis komoditas yang diizinkan
CREATE TABLE komoditas (
<<<<<<< HEAD
    id       VARCHAR(20) PRIMARY KEY,            -- format custom: kmdC-001, kmdT-001
=======
    id       INT AUTO_INCREMENT PRIMARY KEY,
>>>>>>> 47e513f6aa240dc3b450bf0f044fb0c8851a92af
    nama     VARCHAR(50) NOT NULL,               -- beras, jagung, cabai, dll
    satuan   VARCHAR(20) NOT NULL,               -- kg, ton, ikat
    kategori VARCHAR(50)
);

-- products: listing produk dari petani
CREATE TABLE products (
    id           VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id      VARCHAR(36),
<<<<<<< HEAD
    komoditas_id VARCHAR(20),
=======
    komoditas_id INT,
>>>>>>> 47e513f6aa240dc3b450bf0f044fb0c8851a92af
    harga        DECIMAL(12,2) NOT NULL,
    stok         DECIMAL(10,2) NOT NULL,
    foto_url     TEXT,                           -- URL ke GCP Cloud Storage
    status       VARCHAR(20) DEFAULT 'pending'   -- pending / approved / rejected
                 CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (komoditas_id) REFERENCES komoditas(id) ON DELETE CASCADE
);

-- orders: transaksi pembelian (header)
CREATE TABLE orders (
<<<<<<< HEAD
    id                VARCHAR(20) PRIMARY KEY,         -- format custom: ORD-001
=======
    id                VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
>>>>>>> 47e513f6aa240dc3b450bf0f044fb0c8851a92af
    buyer_id          VARCHAR(36),
    total_harga       DECIMAL(14,2) NOT NULL,
    status            VARCHAR(20) DEFAULT 'pending'    -- pending / confirmed / rejected / paid / shipped / done
                      CHECK (status IN ('pending', 'confirmed', 'rejected', 'paid', 'shipped', 'done')),
    payment_token     VARCHAR(255),
    payment_method    VARCHAR(50),
    payment_reference VARCHAR(100),
    paid_at           TIMESTAMP NULL DEFAULT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE CASCADE
);

-- order_items: rincian transaksi pembelian (detail)
CREATE TABLE order_items (
    id         VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
<<<<<<< HEAD
    order_id   VARCHAR(20),
=======
    order_id   VARCHAR(36),
>>>>>>> 47e513f6aa240dc3b450bf0f044fb0c8851a92af
    product_id VARCHAR(36),
    jumlah     DECIMAL(10,2) NOT NULL,
    harga_unit DECIMAL(12,2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- harga_pasar: data historis harga komoditas, ini yang dipakai Prophet
-- di-seed dari BPS atau hargapangan.id, bisa di-update oleh admin
CREATE TABLE harga_pasar (
    id           INT AUTO_INCREMENT PRIMARY KEY,
<<<<<<< HEAD
    komoditas_id VARCHAR(20),
=======
    komoditas_id INT,
>>>>>>> 47e513f6aa240dc3b450bf0f044fb0c8851a92af
    harga        DECIMAL(12,2) NOT NULL,
    wilayah      VARCHAR(100) NOT NULL,
    tanggal      DATE NOT NULL,
    UNIQUE KEY unique_harga_pasar (komoditas_id, wilayah, tanggal),
    FOREIGN KEY (komoditas_id) REFERENCES komoditas(id) ON DELETE CASCADE
);
