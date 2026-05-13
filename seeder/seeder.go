package seeder

import "gorm.io/gorm"

// RunAllSeeders adalah fungsi publik penjalan semua seeder
func RunAllSeeders(db *gorm.DB) {
	seedRole(db)
	seedAdmin(db)
	SeedKaryawan(db)
	SeedPelanggan(db)
}