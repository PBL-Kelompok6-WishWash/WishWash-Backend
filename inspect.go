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

	var order model.Order
	err = db.Preload("RiwayatStatusDetail.ReferensiStatus").First(&order, 33).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Order ID: %d | Kode: %s | Qty: %.2f\n", order.IDOrder, order.KodeOrder, order.Kuantitas)
	fmt.Println("History entries:")
	for i, rs := range order.RiwayatStatusDetail {
		fmt.Printf("  [%d] ID: %d | RefStatusID: %d | Status: %s | Waktu: %s\n", 
			i, rs.IDRiwayat, rs.ReferensiStatusID, rs.ReferensiStatus.NamaStatus, rs.WaktuUpdate)
	}
}
