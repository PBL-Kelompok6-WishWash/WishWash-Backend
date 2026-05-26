# 📚 Dokumentasi Lengkap API WishWash (Postman Collection Reference)

Dokumentasi ini disusun secara lengkap dan mendalam untuk kebutuhan presentasi proyek **WishWash**. API ini berjalan pada **`http://localhost:8080/api/v1`**, terstruktur secara modular berdasarkan peran pengguna (**Admin**, **Karyawan**, **Pelanggan**), dan diamankan dengan skema **Bearer JSON Web Token (JWT)**.

---

## 🗂️ Daftar Isi
1. [🔒 Skema Keamanan & Autentikasi](#-skema-keamanan--autentikasi)
2. [🖥️ API Kelompok Admin](#-api-kelompok-admin)
   - [Manajemen Pelanggan](#manajemen-pelanggan)
   - [Manajemen Karyawan](#manajemen-karyawan)
   - [Manajemen Layanan](#manajemen-layanan)
   - [Manajemen Parfum](#manajemen-parfum)
   - [Manajemen Promo](#manajemen-promo)
   - [Manajemen Metode Pembayaran](#manajemen-metode-pembayaran)
   - [Autentikasi & Profil Admin](#autentikasi--profil-admin)
3. [👷 API Kelompok Karyawan](#-api-kelompok-karyawan)
   - [Autentikasi & Profil Karyawan](#autentikasi--profil-karyawan)
4. [🧑 API Kelompok Pelanggan](#-api-kelompok-pelanggan)
   - [Autentikasi & Profil Pelanggan](#autentikasi--profil-pelanggan)
   - [Manajemen Alamat](#manajemen-alamat)
   - [Katalog Layanan](#katalog-layanan)

---

## 🔒 Skema Keamanan & Autentikasi
Untuk seluruh endpoint yang bertanda `[🔒 JWT SECURE]`, aplikasi wajib menyertakan token JWT pada HTTP Header dengan format sebagai berikut:

* **Header Key**: `Authorization`
* **Header Value**: `Bearer <token_jwt_anda>`

---


---

## 🛑 Standar Respon Gagal (Error) Global
Sebagian besar endpoint dalam API ini mematuhi standar respon error berikut. 
Setiap error akan mengembalikan HTTP Status Code yang relevan beserta format JSON: `{"error": "Deskripsi error"}`.

*   **400 Bad Request**: Dikembalikan ketika format input JSON salah, tipe data tidak sesuai, atau validasi gagal (contoh: ID tidak valid).
*   **401 Unauthorized**: Dikembalikan saat token JWT tidak ada, salah, atau sudah kadaluarsa (khusus endpoint `[🔒 JWT SECURE]`).
*   **403 Forbidden**: Dikembalikan saat user mencoba mengakses endpoint milik *Role* lain (contoh: Karyawan mengakses rute Admin).
*   **404 Not Found**: Dikembalikan ketika data yang dicari (berdasarkan parameter `:id`) tidak ditemukan di database.
*   **500 Internal Server Error**: Dikembalikan saat terjadi kesalahan fatal pada sistem backend atau query database (GORM).


## 🖥️ API Kelompok Admin
Seluruh endpoint di bawah ini dilindungi oleh middleware ganda: `JWTAuthMiddleware` dan `AdminOnly` (Role ID = 1).

### 👥 Manajemen Pelanggan

#### 1. POST Tambah Pelanggan
* **Endpoint**: `POST /admin/pelanggan`
* **Deskripsi**: Mendaftarkan akun pelanggan baru secara manual melalui dashboard admin.
* **Request Body (JSON)**:
  ```json
  {
    "username": "budi_laundry",
    "password": "password123",
    "email": "budi@example.com",
    "nama_lengkap": "Budi Santoso",
    "no_telp": "081234567890",
    "id_role": 3
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Registrasi akun berhasil!",
    "username": "budi_laundry"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **409 Conflict** (Duplikat Data):
    ```json
    {
      "error": "Data duplikat (misal: Username/Email sudah terdaftar)"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. PUT Edit Data Pelanggan
* **Endpoint**: `PUT /admin/pelanggan/:id`
* **Deskripsi**: Mengubah data profil pelanggan berdasarkan ID Pelanggan.
* **Request Body (JSON)**:
  ```json
  {
    "nama_lengkap": "Budi Santoso Nugroho",
    "no_telp": "089876543210"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Data pelanggan berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. DEL Hapus Data Pelanggan
* **Endpoint**: `DELETE /admin/pelanggan/:id`
* **Deskripsi**: Menghapus akun pelanggan secara permanen dari sistem.
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Data pelanggan berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. GET Ambil Data Pelanggan
* **Endpoint**: `GET /admin/pelanggan`
* **Deskripsi**: Mengambil daftar seluruh pelanggan terdaftar di sistem lengkap dengan data relasi akun (`User`) dan perannya (`Role`).
* **Respons Sukses (200 OK - Lengkap Relasi)**:
  ```json
  [
    {
      "id_pelanggan": 1,
      "id_user": 3,
      "nama_lengkap": "Budi Santoso",
      "no_telp": "081234567890",
      "foto_pelanggan": "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde",
      "User": {
        "id_user": 3,
        "id_role": 3,
        "username": "budi_laundry",
        "email": "budi@example.com",
        "Role": {
          "id_role": 3,
          "nama_role": "Pelanggan"
        }
      }
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 💼 Manajemen Karyawan

#### 1. POST Tambah Karyawan
* **Endpoint**: `POST /admin/karyawan`
* **Deskripsi**: Mendaftarkan karyawan operasional laundry baru.
* **Request Body (JSON)**:
  ```json
  {
    "username": "tono_wash",
    "password": "securepassword",
    "email": "tono@example.com",
    "nama_lengkap": "Tono Wijaya",
    "no_telp": "085678901234",
    "id_role": 2
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Registrasi akun berhasil!",
    "username": "tono_wash"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **409 Conflict** (Duplikat Data):
    ```json
    {
      "error": "Data duplikat (misal: Username/Email sudah terdaftar)"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. PUT Edit Data Karyawan
* **Endpoint**: `PUT /admin/karyawan/:id`
* **Deskripsi**: Mengubah profil detail karyawan operasional.
* **Request Body (JSON)**:
  ```json
  {
    "nama_lengkap": "Tono Wijaya Saputra",
    "no_telp": "085699998888",
    "plat_nomor": "BP 1234 XY",
    "status_ketersediaan": "Tersedia"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Data karyawan berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. DEL Hapus Data Karyawan
* **Endpoint**: `DELETE /admin/karyawan/:id`
* **Deskripsi**: Menghapus data karyawan secara permanen dari sistem database.
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Data karyawan berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. GET Ambil Data Karyawan
* **Endpoint**: `GET /admin/karyawan`
* **Deskripsi**: Mengambil daftar seluruh karyawan lengkap dengan plat nomor, status operasional, serta objek relasi akun (`User`) dan perannya (`Role`).
* **Respons Sukses (200 OK - Lengkap Relasi)**:
  ```json
  [
    {
      "id_karyawan": 1,
      "id_user": 4,
      "nama_karyawan": "Tono Wijaya",
      "foto_karyawan": "https://api.dicebear.com/8.x/avataaars-neutral/svg?seed=tono",
      "no_telp": "085678901234",
      "plat_nomor": "BP 1234 XY",
      "jenis_kendaraan": "Motor",
      "status_ketersediaan": "Tersedia",
      "User": {
        "id_user": 4,
        "id_role": 2,
        "username": "tono_wash",
        "email": "tono@example.com",
        "Role": {
          "id_role": 2,
          "nama_role": "Karyawan"
        }
      }
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 🧺 Manajemen Layanan

#### 1. POST Tambah Layanan
* **Endpoint**: `POST /admin/layanan`
* **Deskripsi**: Menambahkan jenis kategori layanan laundry baru.
* **Request Body (JSON)**:
  ```json
  {
    "nama_layanan": "Cuci Kering Lipat",
    "deskripsi_layanan": "Layanan cuci bersih dan dikeringkan otomatis.",
    "harga_per_satuan": 6000,
    "jenis_satuan": "Kg",
    "warna_layanan": "#00BCD4",
    "gambar_layanan": "assets/images/cuci_kering_lipat.png"
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Layanan berhasil ditambahkan",
    "id_layanan": 1
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. PUT Edit Data Layanan
* **Endpoint**: `PUT /admin/layanan/:id`
* **Deskripsi**: Memperbarui harga, deskripsi, warna, atau gambar layanan.
* **Request Body (JSON)**:
  ```json
  {
    "nama_layanan": "Cuci Kering Lipat Super",
    "harga_per_satuan": 6500,
    "warna_layanan": "#00838F"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Layanan berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. DEL Hapus Data Layanan
* **Endpoint**: `DELETE /admin/layanan/:id`
* **Deskripsi**: Menghapus salah satu layanan laundry dari katalog master.
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Layanan berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. GET Ambil Data Layanan
* **Endpoint**: `GET /admin/layanan`
* **Deskripsi**: Mengambil daftar seluruh layanan laundry lengkap dengan **status keaktifan**, tema warna, deskripsi, serta seluruh **referensi tahapan status pesanan** dan **paket durasi layanan** yang tersedia (GORM Preload).
* **Respons Sukses (200 OK - Lengkap Struktur GORM)**:
  ```json
  [
    {
      "id_layanan": 1,
      "nama_layanan": "Cuci Kering Lipat",
      "gambar_layanan": "assets/images/cuci_kering_lipat.png",
      "jenis_satuan": "Kg",
      "harga_per_satuan": 6000,
      "status_layanan": "Aktif",
      "warna_layanan": "#00BCD4",
      "deskripsi_layanan": "Layanan cuci bersih dan dikeringkan otomatis.",
      "referensi_status": [
        {
          "id_referensi_status_layanan": 1,
          "id_layanan": 1,
          "nama_status": "Menunggu Penjemputan",
          "urutan_tahap": 1
        },
        {
          "id_referensi_status_layanan": 2,
          "id_layanan": 1,
          "nama_status": "Pakaian Dijemput",
          "urutan_tahap": 2
        },
        {
          "id_referensi_status_layanan": 3,
          "id_layanan": 1,
          "nama_status": "Sedang Ditimbang",
          "urutan_tahap": 3
        },
        {
          "id_referensi_status_layanan": 4,
          "id_layanan": 1,
          "nama_status": "Proses Cuci",
          "urutan_tahap": 4
        },
        {
          "id_referensi_status_layanan": 5,
          "id_layanan": 1,
          "nama_status": "Siap Diantar",
          "urutan_tahap": 5
        },
        {
          "id_referensi_status_layanan": 6,
          "id_layanan": 1,
          "nama_status": "Selesai",
          "urutan_tahap": 6
        }
      ],
      "paket_layanan": [
        {
          "id_paket_layanan": 4,
          "id_layanan": 1,
          "nama_paket": "Standard",
          "durasi_jam": 72,
          "biaya_tambahan": 0
        },
        {
          "id_paket_layanan": 5,
          "id_layanan": 1,
          "nama_paket": "Premium",
          "durasi_jam": 24,
          "biaya_tambahan": 5000
        },
        {
          "id_paket_layanan": 6,
          "id_layanan": 1,
          "nama_paket": "Express",
          "durasi_jam": 6,
          "biaya_tambahan": 10000
        }
      ]
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 🌹 Manajemen Parfum

#### 1. POST Tambah Parfum
* **Endpoint**: `POST /admin/parfum`
* **Deskripsi**: Menambahkan varian aroma pewangi setrika uap baru.
* **Request Body (JSON)**:
  ```json
  {
    "nama_parfum": "Lavender Bliss",
    "deskripsi_parfum": "Aroma bunga lavender segar penenang pikiran.",
    "status_parfum": "Tersedia"
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Parfum berhasil ditambahkan",
    "id_parfum": 1
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. PUT Edit Data Parfum
* **Endpoint**: `PUT /admin/parfum/:id`
* **Deskripsi**: Memperbarui status ketersediaan aroma parfum.
* **Request Body (JSON)**:
  ```json
  {
    "status_parfum": "Habis"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Parfum berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. DEL Hapus Data Parfum
* **Endpoint**: `DELETE /admin/parfum/:id`
* **Deskripsi**: Menghapus varian parfum dari database.
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Parfum berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. GET Ambil Data Parfum
* **Endpoint**: `GET /admin/parfum`
* **Respons Sukses (200 OK)**:
  ```json
  [
    {
      "id_parfum": 1,
      "nama_parfum": "Lavender Bliss",
      "deskripsi_parfum": "Aroma bunga lavender segar.",
      "status_parfum": "Tersedia"
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 🏷️ Manajemen Promo

#### 1. POST Tambah Promo
* **Endpoint**: `POST /admin/promo`
* **Request Body (JSON)**:
  ```json
  {
    "nama_promo": "Diskon Akhir Bulan",
    "kode_promo": "MEGADEAL",
    "nominal_potongan": 15000,
    "min_order": 50000,
    "tgl_mulai": "2026-05-19T00:00:00Z",
    "tgl_berakhir": "2026-05-31T23:59:59Z",
    "gambar_promo": "https://api.wishwash.my.id/assets/promo/megadeal.png"
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Promo berhasil ditambahkan"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. PUT Edit Data Promo
* **Endpoint**: `PUT /admin/promo/:id`
* **Request Body (JSON)**:
  ```json
  {
    "nominal_potongan": 17000
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Promo berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. DEL Hapus Data Promo
* **Endpoint**: `DELETE /admin/promo/:id`
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Promo berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. GET Ambil Data Promo
* **Endpoint**: `GET /admin/promo`
* **Respons Sukses (200 OK)**:
  ```json
  [
    {
      "id_promo": 1,
      "nama_promo": "Diskon Akhir Bulan",
      "kode_promo": "MEGADEAL",
      "nominal_potongan": 15000,
      "tgl_berakhir": "2026-05-31T23:59:59Z"
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 💳 Manajemen Metode Pembayaran

#### 1. POST Tambah Metode Pembayaran
* **Endpoint**: `POST /admin/metode-pembayaran`
* **Request Body (JSON)**:
  ```json
  {
    "nama_metode": "GoPay QRIS",
    "tipe_metode": "Midtrans",
    "biaya_admin": 0,
    "gambar_metode": "assets/images/qris.png"
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Metode pembayaran berhasil ditambahkan"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. PUT Edit Data Metode Pembayaran
* **Endpoint**: `PUT /admin/metode-pembayaran/:id`
* **Request Body (JSON)**:
  ```json
  {
    "biaya_admin": 1000
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Metode pembayaran berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. DEL Hapus Data Metode Pembayaran
* **Endpoint**: `DELETE /admin/metode-pembayaran/:id`
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Metode pembayaran berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. GET Ambil Data Metode Pembayaran
* **Endpoint**: `GET /admin/metode-pembayaran`
* **Respons Sukses (200 OK)**:
  ```json
  [
    {
      "id_metode_pembayaran": 1,
      "nama_metode": "GoPay QRIS",
      "tipe_metode": "Midtrans",
      "biaya_admin": 0,
      "gambar_metode": "assets/images/qris.png"
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 🔑 Autentikasi & Profil Admin

#### 1. POST Login Admin
* **Endpoint**: `POST /auth/login`
* **Request Body (JSON)**:
  ```json
  {
    "username": "superadmin",
    "password": "adminpassword"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Autentikasi berhasil.",
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "id_role": 1,
    "display_name": "Mega Admin"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. GET Profile Admin `[🔒 JWT SECURE]`
* **Endpoint**: `GET /profile`
* **Respons Sukses (200 OK)**:
  ```json
  {
    "id_user": 1,
    "username": "superadmin",
    "email": "admin@wishwash.com",
    "nama_lengkap": "Mega Admin",
    "role": "Admin"
  }
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. PUT Edit Profile Admin `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /profile/update`
* **Request Body (JSON)**:
  ```json
  {
    "nama_lengkap": "Mega Admin WishWash"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Profil berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. PUT Edit Password Admin `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /profile/password`
* **Request Body (JSON)**:
  ```json
  {
    "password_lama": "adminpassword",
    "password_baru": "newadminpassword"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Kata sandi berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

## 👷 API Kelompok Karyawan
Seluruh endpoint di bawah ini dilindungi oleh middleware ganda: `JWTAuthMiddleware` dan `KaryawanOrAdmin` (Role ID = 2).

### 🔑 Autentikasi & Profil Karyawan

#### 1. POST Login Karyawan
* **Endpoint**: `POST /auth/login`
* **Request Body (JSON)**:
  ```json
  {
    "username": "tono_wash",
    "password": "securepassword"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Autentikasi berhasil.",
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "id_role": 2,
    "display_name": "Tono Wijaya"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **409 Conflict** (Duplikat Data):
    ```json
    {
      "error": "Data duplikat (misal: Username/Email sudah terdaftar)"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. GET Profile Karyawan `[🔒 JWT SECURE]`
* **Endpoint**: `GET /profile`
* **Respons Sukses (200 OK)**:
  ```json
  {
    "id_user": 4,
    "username": "tono_wash",
    "email": "tono@example.com",
    "nama_lengkap": "Tono Wijaya",
    "no_telp": "085678901234",
    "plat_nomor": "BP 1234 XY",
    "status_ketersediaan": "Tersedia",
    "role": "Karyawan"
  }
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. PUT Edit Profile Karyawan `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /profile/update`
* **Request Body (JSON)**:
  ```json
  {
    "plat_nomor": "BP 4321 YX",
    "status_ketersediaan": "Sibuk"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Profil berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. PUT Edit Password Karyawan `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /profile/password`
* **Request Body (JSON)**:
  ```json
  {
    "password_lama": "securepassword",
    "password_baru": "newsecurepassword"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Kata sandi berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

## 🧑 API Kelompok Pelanggan

### 🔑 Autentikasi & Profil Pelanggan

#### 1. POST Login Pelanggan
* **Endpoint**: `POST /auth/login`
* **Request Body (JSON)**:
  ```json
  {
    "username": "budi_laundry",
    "password": "password123"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "id_role": 3,
    "display_name": "Budi Santoso"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **409 Conflict** (Duplikat Data):
    ```json
    {
      "error": "Data duplikat (misal: Username/Email sudah terdaftar)"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. POST Register Pelanggan
* **Endpoint**: `POST /auth/register`
* **Request Body (JSON)**:
  ```json
  {
    "username": "cindy_cute",
    "password": "cindypassword",
    "email": "cindy@example.com",
    "nama_lengkap": "Cindy Aulia",
    "no_telp": "087711223344",
    "id_role": 3
  }
  ```
* **Respons Sukses (201 Created)**:
  ```json
  {
    "message": "Registrasi akun berhasil!",
    "username": "cindy_cute"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **409 Conflict** (Duplikat Data):
    ```json
    {
      "error": "Data duplikat (misal: Username/Email sudah terdaftar)"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. GET Profile Pelanggan `[🔒 JWT SECURE]`
* **Endpoint**: `GET /profile`
* **Respons Sukses (200 OK)**:
  ```json
  {
    "id_user": 3,
    "username": "budi_laundry",
    "email": "budi@example.com",
    "nama_lengkap": "Budi Santoso",
    "no_telp": "081234567890",
    "role": "Pelanggan"
  }
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. PUT Edit Password Pelanggan `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /profile/password`
* **Request Body (JSON)**:
  ```json
  {
    "password_lama": "password123",
    "password_baru": "newpassword123"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Kata sandi berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 🏠 Manajemen Alamat

#### 1. GET Ambil Daftar Alamat `[🔒 JWT SECURE]`
* **Endpoint**: `GET /alamat`
* **Deskripsi**: Mengambil seluruh alamat yang didaftarkan oleh pelanggan aktif.
* **Respons Sukses (200 OK)**:
  ```json
  [
    {
      "id_alamat": 1,
      "label_alamat": "Rumah Utama",
      "alamat_lengkap": "Jl. Sudirman No. 45, Batam Center",
      "is_primary": true
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 2. POST Tambah Alamat Baru `[🔒 JWT SECURE]`
* **Endpoint**: `POST /alamat`
* **Request Body (JSON)**:
  ```json
  {
    "label_alamat": "Kantor Cabang",
    "alamat_lengkap": "Komp. Nagoya Hill Blok C No. 12",
    "is_primary": false
  }
  ```
* **Respons Sukses (210 Created)**:
  ```json
  {
    "message": "Alamat baru berhasil ditambahkan",
    "id_alamat": 2
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 3. PUT Edit Alamat `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /alamat/:id`
* **Request Body (JSON)**:
  ```json
  {
    "label_alamat": "Kos Baru",
    "alamat_lengkap": "Jl. Kepri Raya No. 10"
  }
  ```
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Alamat berhasil diperbarui"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 4. PUT Set Alamat Menjadi Utama `[🔒 JWT SECURE]`
* **Endpoint**: `PUT /alamat/:id/primary`
* **Deskripsi**: Menjadikan alamat yang dipilih sebagai alamat default utama penjemputan/penyerahan laundry.
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Alamat utama berhasil disetel"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (Validasi Gagal/Format JSON Salah):
    ```json
    {
      "error": "Format data tidak valid atau field wajib kosong"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


#### 5. DEL Hapus Alamat `[🔒 JWT SECURE]`
* **Endpoint**: `DELETE /alamat/:id`
* **Respons Sukses (200 OK)**:
  ```json
  {
    "message": "Alamat berhasil dihapus"
  }
  ```

* **Respons Error (Gagal)**:
  * **400 Bad Request** (ID Tidak Valid):
    ```json
    {
      "error": "ID parameter tidak valid"
    }
    ```
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **404 Not Found** (Data Kosong):
    ```json
    {
      "error": "Data tidak ditemukan di database"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```


---

### 📖 Katalog Layanan

#### 1. GET Ambil Daftar Layanan `[🔒 JWT SECURE]`
* **Endpoint**: `GET /layanan`
* **Deskripsi**: Mengambil daftar layanan laundry terpadu yang aktif bagi Pelanggan lengkap dengan seluruh referensi tahapan status pesanan dan paket durasi layanan (GORM Preload).
* **Respons Sukses (200 OK - Lengkap Struktur GORM)**:
  ```json
  [
    {
      "id_layanan": 1,
      "nama_layanan": "Cuci Kering Lipat",
      "gambar_layanan": "assets/images/cuci_kering_lipat.png",
      "jenis_satuan": "Kg",
      "harga_per_satuan": 6000,
      "status_layanan": "Aktif",
      "warna_layanan": "#00BCD4",
      "deskripsi_layanan": "Layanan cuci bersih dan dikeringkan otomatis.",
      "referensi_status": [
        {
          "id_referensi_status_layanan": 1,
          "id_layanan": 1,
          "nama_status": "Menunggu Penjemputan",
          "urutan_tahap": 1
        },
        {
          "id_referensi_status_layanan": 2,
          "id_layanan": 1,
          "nama_status": "Pakaian Dijemput",
          "urutan_tahap": 2
        },
        {
          "id_referensi_status_layanan": 3,
          "id_layanan": 1,
          "nama_status": "Sedang Ditimbang",
          "urutan_tahap": 3
        },
        {
          "id_referensi_status_layanan": 4,
          "id_layanan": 1,
          "nama_status": "Proses Cuci",
          "urutan_tahap": 4
        },
        {
          "id_referensi_status_layanan": 5,
          "id_layanan": 1,
          "nama_status": "Siap Diantar",
          "urutan_tahap": 5
        },
        {
          "id_referensi_status_layanan": 6,
          "id_layanan": 1,
          "nama_status": "Selesai",
          "urutan_tahap": 6
        }
      ],
      "paket_layanan": [
        {
          "id_paket_layanan": 4,
          "id_layanan": 1,
          "nama_paket": "Standard",
          "durasi_jam": 72,
          "biaya_tambahan": 0
        },
        {
          "id_paket_layanan": 5,
          "id_layanan": 1,
          "nama_paket": "Premium",
          "durasi_jam": 24,
          "biaya_tambahan": 5000
        },
        {
          "id_paket_layanan": 6,
          "id_layanan": 1,
          "nama_paket": "Express",
          "durasi_jam": 6,
          "biaya_tambahan": 10000
        }
      ]
    }
  ]
  ```

* **Respons Error (Gagal)**:
  * **401 Unauthorized** (Token Kosong/Kadaluarsa):
    ```json
    {
      "error": "Token JWT tidak valid atau tidak ditemukan"
    }
    ```
  * **403 Forbidden** (Hak Akses Ditolak):
    ```json
    {
      "error": "Akses ditolak: Anda tidak memiliki peran yang sesuai"
    }
    ```
  * **500 Internal Server Error** (Gangguan Sistem):
    ```json
    {
      "error": "Gagal memproses data atau kesalahan internal server"
    }
    ```

