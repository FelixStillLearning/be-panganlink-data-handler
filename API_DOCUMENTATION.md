# PanganLink Data Handler API Documentation

Berikut adalah dokumentasi API yang disediakan oleh Data Handler Service PanganLink.

## Base URL
`/api/v1`

## Authentication
Sebagian besar endpoint memerlukan Bearer Token (`Authorization: Bearer <token>`).

## Endpoints

### 1. Auth
- `POST /auth/register`: Mendaftarkan pengguna baru
- `POST /auth/login`: Login pengguna
- `GET /auth/me`: Mendapatkan data pengguna saat ini

### 2. Admin
- `GET /admin/dashboard`: Data statistik admin
- `GET /admin/users`: Daftar pengguna
- `GET /admin/products`: Daftar produk untuk moderasi
- `GET /admin/commodities`: Daftar komoditas utama

### 3. Petani
- `GET /petani/products`: Daftar produk petani
- `POST /petani/products`: Tambah produk baru
- `GET /petani/orders`: Daftar pesanan masuk
- `PUT /petani/orders/:id/status`: Update status pesanan

### 4. Pembeli
- `GET /pembeli/products`: Daftar produk yang tersedia
- `POST /pembeli/orders`: Membuat pesanan (checkout)
- `GET /pembeli/history`: Riwayat pembelian
