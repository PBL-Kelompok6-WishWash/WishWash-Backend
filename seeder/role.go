package seeder // <-- Perhatikan, packagenya ganti jadi seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

func seedRole(db *gorm.DB) {
	var count int64
	db.Model(&model.Role{}).Count(&count)

	if count == 0 {
		roles := []model.Role{
			{IDRole: 1, NamaRole: "Admin"},
			{IDRole: 2, NamaRole: "Karyawan"},
			{IDRole: 3, NamaRole: "Pelanggan"},
		}

		if err := db.Create(&roles).Error; err != nil {
			log.Println("❌ Gagal seeding tabel Role:", err)
		} else {
			log.Println("🌱 Seeding Role berhasil!")
		}
	} else {
		log.Println("✅ Tabel Role sudah terisi.")
	}
}