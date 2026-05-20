package seeder

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

func SeedLayanan(db *gorm.DB) {
	layanans := []model.Layanan{
		{
			NamaLayanan:    "Cuci Kering Lipat",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 7000,
			WarnaLayanan:   "#00BCD4", // Cyan
			GambarLayanan:  "",
			DeskripsiLayanan: "Paket lengkap cuci bersih, kering, lipat rapi.",
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Penjemputan", UrutanTahap: 2},
				{NamaStatus: "Proses Timbang", UrutanTahap: 3},
				{NamaStatus: "Proses Cuci", UrutanTahap: 4},
				{NamaStatus: "Proses Kering", UrutanTahap: 5},
				{NamaStatus: "Proses Lipat", UrutanTahap: 6},
				{NamaStatus: "Siap Diantar", UrutanTahap: 7},
				{NamaStatus: "Selesai", UrutanTahap: 8},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 48, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 24, BiayaTambahan: 5000},
				{NamaPaket: "Kilat", DurasiJam: 6, BiayaTambahan: 10000},
			},
			StatusLayanan: "Aktif",
		},
		{
			NamaLayanan:    "Cuci Kering",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 5000,
			WarnaLayanan:   "#8BC34A", // Hijau muda
			GambarLayanan:  "",
			DeskripsiLayanan: "Dicuci bersih dan dikeringkan tanpa disetrika.",
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Penjemputan", UrutanTahap: 2},
				{NamaStatus: "Proses Timbang", UrutanTahap: 3},
				{NamaStatus: "Proses Cuci", UrutanTahap: 4},
				{NamaStatus: "Proses Kering", UrutanTahap: 5},
				{NamaStatus: "Siap Diantar", UrutanTahap: 6},
				{NamaStatus: "Selesai", UrutanTahap: 7},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 48, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 24, BiayaTambahan: 3000},
				{NamaPaket: "Kilat", DurasiJam: 6, BiayaTambahan: 7000},
			},
			StatusLayanan: "Aktif",
		},
		{
			NamaLayanan:    "Cuci & Setrika",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 10000,
			WarnaLayanan:   "#9C27B0", // Ungu
			GambarLayanan:  "",
			DeskripsiLayanan: "Dicuci bersih, wangi, dan disetrika rapi.",
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Penjemputan", UrutanTahap: 2},
				{NamaStatus: "Proses Timbang", UrutanTahap: 3},
				{NamaStatus: "Proses Cuci", UrutanTahap: 4},
				{NamaStatus: "Proses Kering", UrutanTahap: 5},
				{NamaStatus: "Proses Setrika", UrutanTahap: 6},
				{NamaStatus: "Siap Diantar", UrutanTahap: 7},
				{NamaStatus: "Selesai", UrutanTahap: 8},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 48, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 24, BiayaTambahan: 6000},
				{NamaPaket: "Kilat", DurasiJam: 6, BiayaTambahan: 12000},
			},
			StatusLayanan: "Aktif",
		},
		{
			NamaLayanan:    "Setrika",
			JenisSatuan:    "Kg",
			HargaPerSatuan: 6000,
			WarnaLayanan:   "#FFC107", // Kuning
			GambarLayanan:  "",
			DeskripsiLayanan: "Pakaian disetrika rapi & harum premium.",
			ReferensiStatus: []model.ReferensiStatusLayanan{
				{NamaStatus: "Pesanan Diterima", UrutanTahap: 1},
				{NamaStatus: "Penjemputan", UrutanTahap: 2},
				{NamaStatus: "Proses Timbang", UrutanTahap: 3},
				{NamaStatus: "Proses Setrika", UrutanTahap: 4},
				{NamaStatus: "Siap Diantar", UrutanTahap: 5},
				{NamaStatus: "Selesai", UrutanTahap: 6},
			},
			PaketLayanan: []model.PaketLayanan{
				{NamaPaket: "Reguler", DurasiJam: 24, BiayaTambahan: 0},
				{NamaPaket: "Express", DurasiJam: 12, BiayaTambahan: 4000},
				{NamaPaket: "Kilat", DurasiJam: 4, BiayaTambahan: 8000},
			},
			StatusLayanan: "Aktif",
		},
	}

	for _, l := range layanans {
		var existing model.Layanan
		// Cek berdasarkan NamaLayanan
		if err := db.Where("nama_layanan = ?", l.NamaLayanan).First(&existing).Error; err == nil {
			log.Printf("ℹ️ Layanan %s sudah ada, dilewati (seeder tidak meng-update data).\n", l.NamaLayanan)
		} else {
			// Buat baru jika belum ada
			if err := db.Create(&l).Error; err != nil {
				log.Printf("❌ Gagal seeding layanan %s: %v\n", l.NamaLayanan, err)
			} else {
				log.Printf("🌱 Berhasil seeding layanan %s!\n", l.NamaLayanan)
			}
		}
	}
}
