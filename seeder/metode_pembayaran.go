package seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

func SeedMetodePembayaran(db *gorm.DB) {
	var count int64
	db.Model(&model.MetodePembayaran{}).Count(&count)

	if count > 0 {
		log.Println("✅ Data Metode Pembayaran sudah ada, skip proses seeding.")
		return
	}

	mps := []model.MetodePembayaran{
		{
			NamaMetode:   "Tunai",
			TipeMetode:   "Tunai",
			KodeMetode:   "cash",
			StatusMetode: "Aktif",
			GambarMetode: "",
		},
		{
			NamaMetode:   "QRIS",
			TipeMetode:   "Midtrans",
			KodeMetode:   "qris",
			StatusMetode: "Aktif",
			GambarMetode: "",
		},
	}

	for _, mp := range mps {
		if err := db.Create(&mp).Error; err != nil {
			log.Printf("❌ Gagal seeding metode pembayaran %s: %v\n", mp.NamaMetode, err)
		} else {
			log.Printf("🌱 Berhasil seeding metode pembayaran %s!\n", mp.NamaMetode)
		}
	}
}
