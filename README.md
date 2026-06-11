<div align="center">
  <h1>⚙️ PanganLink - Data Handler Service</h1>
  <p><em>Layanan backend utama dan gateway untuk manajemen data PanganLink.</em></p>
  
  ![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
  ![MySQL](https://img.shields.io/badge/MySQL-00000F?style=for-the-badge&logo=mysql&logoColor=white)
  ![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
</div>

---

## 📖 Tentang Layanan
Repositori ini berisi *Data Handler Service* yang dikembangkan menggunakan bahasa Go (Golang). Layanan ini bertanggung jawab atas pemrosesan data utama, autentikasi, manajemen komoditas, transaksi, serta menjadi jembatan integrasi bagi *frontend* dan *AI service*.

## 👥 Tim Pengembang
Proyek ini dikembangkan untuk memenuhi tugas mata kuliah Komputasi Awan oleh:

| NRP | Nama |
| :--- | :--- |
| `152023018` | Ghinova Klarisa Irawadi |
| `152023148` | Felix Angga Resky |
| `152023141` | Parisan Apro |
| `152023167` | Raelqiansyah Putranta Dibrata |
| `152023186` | Difie Anggely |

## 📚 Dokumentasi
- [API Documentation](./API_DOCUMENTATION.md) - Detail endpoint *gateway service*.

## 🚀 Setup & Instalasi

### Persyaratan
- **Go** (Versi 1.22+)
- **Docker** (Opsional untuk containerisasi)
- **Database** (MySQL / PostgreSQL)

### Langkah-langkah
1. Salin file environment:
   ```bash
   cp .env.example .env
   ```
2. Konfigurasi kredensial database dan environment variabel di `.env`.
3. Jalankan server secara lokal:
   ```bash
   make run
   ```

## 🛠️ Perintah Makefile
| Perintah | Deskripsi |
| :--- | :--- |
| `make run` | Menjalankan *development server* |
| `make build` | Melakukan *build* file *binary* |
| `make test` | Menjalankan *unit testing* |
| `make lint` | Menjalankan *linter* untuk merapikan kode |
| `make docker-build` | Membangun *Docker image* aplikasi |
