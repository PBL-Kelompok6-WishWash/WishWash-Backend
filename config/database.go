package config

import (
	"fmt"
	"log"
	// "os"

	"github.com/PBL-Kelompok6-WishWash/backend/model" // Import model
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Membaca kredensial dari environment variables (atau hardcode sementara untuk tes)
	host := "localhost"
	user := "postgres"
	password := "12345678" // Password database
	dbname := "wishwash_db"
	port := "5433"         // Pastikan port-nya 5433

	// Konfigurasi string koneksi
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	// Membuka koneksi menggunakan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal terhubung ke database: ", err)
	}

	log.Println("✅ Berhasil terhubung ke database PostgreSQL (Port 5433)!")

	err = database.AutoMigrate(
		// --- Modul Autentikasi & Pengguna ---
		&model.Role{},
		&model.User{},
		&model.Admin{},
		&model.Pelanggan{},
		&model.Karyawan{},

		// --- Modul Master Operasional Laundry ---
		&model.Layanan{},
		&model.PaketLayanan{},
		&model.Parfum{},
		&model.Promo{},
		&model.MetodePembayaran{},

		// --- Modul Pendukung & Konfigurasi ---
		&model.Alamat{},
		&model.Settings{},
		&model.Notifikasi{},
		&model.ReferensiStatusLayanan{},

		// --- Modul Transaksi Inti ---
		&model.Order{},
		&model.PromoOrder{},

		// --- Modul Pasca-Order & Tracking ---
		&model.Pembayaran{},
		&model.Penilaian{},
		&model.TrackingPengiriman{},
		&model.RiwayatStatusDetail{},

		// --- Modul Chat & Komunikasi ---
		&model.RoomChat{},
		&model.PesanChat{},
		&model.ChatGambar{},
	)
	if err != nil {
		log.Fatal("❌ Gagal melakukan migrasi tabel: ", err)
	}
	log.Println("🚀 Migrasi Tabel ke Database berhasil!")

	DB = database
}