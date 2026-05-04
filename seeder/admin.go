package seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedAdmin(db *gorm.DB) {
	var count int64
	// Cek apakah sudah ada admin di tabel User
	db.Model(&model.User{}).Where("id_role = ?", 1).Count(&count)

	if count == 0 {
		// 1. Buat Password Terenkripsi
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

		// 2. Buat Akun User untuk Admin
		adminUser := model.User{
			Username: "admin_wishwash",
			Email:    "admin@wishwash.com",
			Password: string(hashedPassword),
			RoleID:   1, // Pastikan ID 1 adalah Role Admin
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Println("❌ Gagal seeding User Admin:", err)
			return
		}

		// 3. Buat Profil Admin di tabel Admin (Gunakan ID dari adminUser yang baru dibuat)
		// Pastikan struct model.Admin kamu punya field IDUser dan Nama
		
		adminProfile := model.Admin{
			UserID: adminUser.IDUser,
			NamaAdmin:   "Jono",
		}
		
		if err := db.Create(&adminProfile).Error; err != nil {
			log.Println("❌ Gagal seeding Profil Admin:", err)
		} else {
			log.Println("🌱 Seeding Akun & Profil Admin berhasil!")
		}
		
		log.Println("🌱 Seeding Akun Admin berhasil! (Profil Admin menyusul setelah model siap)")
	} else {
		log.Println("✅ Akun Admin sudah tersedia.")
	}
}