package seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedPelanggan(db *gorm.DB) {
	// 1. Cek apakah tabel pelanggan sudah ada isinya biar nggak duplikat pas di-run berkali-kali
	var count int64
	if err := db.Model(&model.Pelanggan{}).Count(&count).Error; err != nil {
		log.Printf("Gagal mengecek data pelanggan: %v", err)
		return
	}

	if count > 0 {
		log.Println("Seeder: Data Pelanggan sudah ada, skip proses seeding.")
		return
	}

	// 2. Siapkan password default (misal: "password123")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// 3. Array data dummy
	dummyPelanggans := []struct {
		Username    string
		Email       string
		NamaLengkap string
		NoTelp      string
		Foto        string
	}{
		{"budi_santoso", "budi@gmail.com", "Budi Santoso", "081234567890", "https://api.dicebear.com/8.x/avataaars/svg?seed=Budi"},
		{"siti_aminah", "siti@yahoo.com", "Siti Aminah", "081987654321", "https://api.dicebear.com/8.x/avataaars/svg?seed=Siti"},
		{"alex_wijaya", "alex@outlook.com", "Alex Wijaya", "082233445566", ""}, // Sengaja kosong biar ngetes fallback ikon
		{"rachel_ven", "rachel@gmail.com", "Rachel Vennya", "087788990011", "https://api.dicebear.com/8.x/avataaars/svg?seed=Rachel"},
		{"dani_kurnia", "dani@gmail.com", "Dani Kurniawan", "085612341234", ""},
	}

	// 4. Looping untuk insert ke Database
	for _, dp := range dummyPelanggans {
		// Insert User Dulu (RoleID 3 untuk Pelanggan)
		user := model.User{
			Username: dp.Username,
			Email:    dp.Email,
			Password: string(hashedPassword),
			RoleID:   3, 
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Seeder: Gagal membuat user untuk %s: %v", dp.Username, err)
			continue
		}

		// Insert Pelanggan pakai ID User yang baru dibuat
		pelanggan := model.Pelanggan{
			UserID:        user.IDUser,
			NamaLengkap:   dp.NamaLengkap,
			NoTelp:        dp.NoTelp,
			FotoPelanggan: dp.Foto,
		}

		if err := db.Create(&pelanggan).Error; err != nil {
			log.Printf("Seeder: Gagal membuat data pelanggan %s: %v", dp.NamaLengkap, err)
		}
	}

	log.Println("Seeder: Berhasil menabur data dummy Pelanggan! 🌱")
}