package seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

func SeedParfum(db *gorm.DB) {
	var count int64
	db.Model(&model.Parfum{}).Count(&count)

	if count > 0 {
		log.Println("✅ Data Parfum sudah ada, skip proses seeding.")
		return
	}

	parfums := []model.Parfum{
		{NamaParfum: "Malaikat Subuh", Keterangan: "Aroma lembut dan menenangkan, cocok untuk pakaian sehari-hari"},
		{NamaParfum: "Lavender Bliss", Keterangan: "Wangi bunga lavender asli yang memberikan ketenangan"},
		{NamaParfum: "Citrus Burst", Keterangan: "Aroma jeruk segar yang memberikan energi dan kesegaran"},
		{NamaParfum: "Fresh Cotton", Keterangan: "Wangi kapas bersih khas cucian baru yang tahan lama"},
		{NamaParfum: "Ocean Breeze", Keterangan: "Aroma laut segar yang ringan dan menyegarkan"},
		{NamaParfum: "Rose Garden", Keterangan: "Wewangian mawar elegan untuk pakaian formal dan spesial"},
	}

	for _, p := range parfums {
		if err := db.Create(&p).Error; err != nil {
			log.Printf("❌ Gagal seeding parfum %s: %v\n", p.NamaParfum, err)
		} else {
			log.Printf("🌱 Berhasil seeding parfum %s!\n", p.NamaParfum)
		}
	}
}
