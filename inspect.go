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

	var notifs []model.Notifikasi
	err = db.Order("id_notifikasi desc").Limit(5).Find(&notifs).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== LAST 5 NOTIFICATIONS ===")
	for _, n := range notifs {
		fmt.Printf("ID: %d | UserID: %d | Title: %s | Message: %s | Created: %s\n", 
			n.IDNotifikasi, n.UserID, n.Judul, n.Pesan, n.CreatedAt)
	}
}
