package seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

func SeedLayanan(db *gorm.DB) {
	var count int64
	db.Model(&model.Layanan{}).Count(&count)

	if count > 0 {
		log.Println("✅ Data Layanan sudah ada, skip proses seeding.")
		return
	}

	layanans := []model.Layanan{
		{
			NamaLayanan:    "Cuci Kering Lipat",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 7000,
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Proses Cuci", UrutanTahap: 2},
				{NamaStatus: "Proses Kering", UrutanTahap: 3},
				{NamaStatus: "Proses Lipat", UrutanTahap: 4},
				{NamaStatus: "Siap Diambil", UrutanTahap: 5},
				{NamaStatus: "Selesai", UrutanTahap: 6},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 48, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 24, BiayaTambahan: 5000},
				{NamaPaket: "Kilat", DurasiJam: 6, BiayaTambahan: 10000},
			},
		},
		{
			NamaLayanan:    "Cuci Kering",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 5000,
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Proses Cuci", UrutanTahap: 2},
				{NamaStatus: "Proses Kering", UrutanTahap: 3},
				{NamaStatus: "Siap Diambil", UrutanTahap: 4},
				{NamaStatus: "Selesai", UrutanTahap: 5},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 48, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 24, BiayaTambahan: 3000},
				{NamaPaket: "Kilat", DurasiJam: 6, BiayaTambahan: 7000},
			},
		},
		{
			NamaLayanan:    "Cuci & Setrika",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 10000,
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Proses Cuci", UrutanTahap: 2},
				{NamaStatus: "Proses Kering", UrutanTahap: 3},
				{NamaStatus: "Proses Setrika", UrutanTahap: 4},
				{NamaStatus: "Siap Diambil", UrutanTahap: 5},
				{NamaStatus: "Selesai", UrutanTahap: 6},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 48, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 24, BiayaTambahan: 6000},
				{NamaPaket: "Kilat", DurasiJam: 6, BiayaTambahan: 12000},
			},
		},
		{
			NamaLayanan:    "Setrika",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 6000,
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Proses Setrika", UrutanTahap: 2},
				{NamaStatus: "Siap Diambil", UrutanTahap: 3},
				{NamaStatus: "Selesai", UrutanTahap: 4},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 24, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 12, BiayaTambahan: 4000},
				{NamaPaket: "Kilat", DurasiJam: 4, BiayaTambahan: 8000},
			},
		},
	}

	for _, l := range layanans {
		if err := db.Create(&l).Error; err != nil {
			log.Printf("❌ Gagal seeding layanan %s: %v\n", l.NamaLayanan, err)
		} else {
			log.Printf("🌱 Berhasil seeding layanan %s!\n", l.NamaLayanan)
		}
	}
}
