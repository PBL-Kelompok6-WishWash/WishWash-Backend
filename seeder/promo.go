package seeder

import (
	"log"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

func SeedPromo(db *gorm.DB) {
	var count int64
	db.Model(&model.Promo{}).Count(&count)

	if count > 0 {
		log.Println("✅ Data Promo sudah ada, skip proses seeding.")
		return
	}

	promos := []model.Promo{
		{
			KodePromo:        "WISHNEW26",
			NamaPromo:        "First Order Promo",
			Deskripsi:        "Diskon khusus pengguna baru WishWash",
			TipePromo:        "Persentase",
			NominalPotongan:  20,
			MinimalOrder:     0,
			MaksimalPotongan: 20000,
			TglMulai:         time.Now(),
			TglBerakhir:      time.Now().AddDate(0, 1, 0),
			StatusPromo:      "Aktif",
			GambarPromo:      "",
		},
		{
			KodePromo:        "WASH5K",
			NamaPromo:        "Payday Flash Wash",
			Deskripsi:        "Potongan langsung di hari gajian",
			TipePromo:        "Nominal",
			NominalPotongan:  5000,
			MinimalOrder:     30000,
			MaksimalPotongan: 5000,
			TglMulai:         time.Now(),
			TglBerakhir:      time.Now().AddDate(0, 0, 7),
			StatusPromo:      "Aktif",
			GambarPromo:      "",
		},
		{
			KodePromo:        "BERSIHHEMAT",
			NamaPromo:        "Promo Bersih Hemat",
			Deskripsi:        "Cuci banyak lebih murah",
			TipePromo:        "Persentase",
			NominalPotongan:  15,
			MinimalOrder:     50000,
			MaksimalPotongan: 15000,
			TglMulai:         time.Now(),
			TglBerakhir:      time.Now().AddDate(0, 2, 0),
			StatusPromo:      "Aktif",
			GambarPromo:      "",
		},
		{
			KodePromo:        "WEEKENDSERU",
			NamaPromo:        "Weekend Seru",
			Deskripsi:        "Laundry santai di akhir pekan",
			TipePromo:        "Nominal",
			NominalPotongan:  10000,
			MinimalOrder:     75000,
			MaksimalPotongan: 10000,
			TglMulai:         time.Now(),
			TglBerakhir:      time.Now().AddDate(0, 3, 0),
			StatusPromo:      "Aktif",
			GambarPromo:      "",
		},
		{
			KodePromo:        "LOYALTY10",
			NamaPromo:        "Loyalty Member",
			Deskripsi:        "Apresiasi untuk pelanggan setia",
			TipePromo:        "Persentase",
			NominalPotongan:  10,
			MinimalOrder:     20000,
			MaksimalPotongan: 10000,
			TglMulai:         time.Now(),
			TglBerakhir:      time.Now().AddDate(1, 0, 0),
			StatusPromo:      "Aktif",
			GambarPromo:      "",
		},
	}

	for _, p := range promos {
		db.Create(&p)
	}

	log.Println("✅ Seeder: 5 Data Promo berhasil ditambahkan!")
}
