package main

import (
	"fmt"
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=12345678 dbname=wishwash_db port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var o model.Order
	err = db.Preload("Layanan.ReferensiStatus").
		Preload("RiwayatStatusDetail.ReferensiStatus").
		Where("kode_order = ?", "WW-E0OWSK").
		First(&o).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order: %s, Layanan: %s\n", o.KodeOrder, o.Layanan.NamaLayanan)
	fmt.Println("  History Status Detail:")
	for _, r := range o.RiwayatStatusDetail {
		fmt.Printf("    - ID: %d, Status: %s, Waktu: %s\n", r.IDRiwayat, r.ReferensiStatus.NamaStatus, r.WaktuUpdate)
	}
	fmt.Println("  Referensi Status Layanan:")
	for _, rs := range o.Layanan.ReferensiStatus {
		fmt.Printf("    - ID: %d, Status: %s, Urutan: %d\n", rs.IDReferensiStatus, rs.NamaStatus, rs.UrutanTahap)
	}
}
