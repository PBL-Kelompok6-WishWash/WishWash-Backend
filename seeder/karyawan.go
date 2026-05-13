package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedKaryawan(db *gorm.DB) {
	var count int64
	if err := db.Model(&model.Karyawan{}).Count(&count).Error; err != nil {
		log.Printf("Gagal mengecek data karyawan: %v", err)
		return
	}

	if count > 0 {
		log.Println("Seeder: Data Karyawan sudah ada, skip proses seeding.")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPasswordSpecific, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

	// 1. Buat Karyawan 1 khusus
	user1 := model.User{
		Username: "karyawan1",
		Email:    "karyawan1@wishwash.com",
		Password: string(hashedPasswordSpecific),
		RoleID:   2,
	}
	db.Create(&user1)
	karyawan1 := model.Karyawan{
		UserID:             user1.IDUser,
		NamaKaryawan:       "Karyawan Spesifik 1",
		NoTelp:             "082222222222",
		PlatNomor:          "BP 1234 AA",
		JenisKendaraan:     "Motor",
		StatusKetersediaan: "Tersedia",
	}
	db.Create(&karyawan1)

	firstNames := []string{"Budi", "Andi", "Roni", "Siti", "Ayu", "Joko", "Tono", "Dewi", "Rina", "Agus", "Bagus", "Cahyo", "Dina", "Eka", "Fajar", "Gita", "Hadi", "Indra", "Jamil", "Kartika"}
	lastNames := []string{"Santoso", "Wijaya", "Kurniawan", "Setiawan", "Pratama", "Putra", "Putri", "Sari", "Lestari", "Hidayat", "Nugroho", "Saputra", "Wahyudi", "Ramadhan"}

	rand.Seed(time.Now().UnixNano())

	// 2. Buat 10 Karyawan Acak
	for i := 1; i <= 10; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		nama := fmt.Sprintf("%s %s", firstName, lastName)
		username := strings.ToLower(fmt.Sprintf("%skaryawan%d", firstName, rand.Intn(100)))
		email := fmt.Sprintf("%s@wishwash.com", username)
		noTelp := fmt.Sprintf("082345%06d", rand.Intn(999999))

		user := model.User{
			Username: username,
			Email:    email,
			Password: string(hashedPassword),
			RoleID:   2,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Seeder: Gagal membuat user untuk %s: %v", username, err)
			continue
		}

		karyawan := model.Karyawan{
			UserID:             user.IDUser,
			NamaKaryawan:       nama,
			NoTelp:             noTelp,
			PlatNomor:          fmt.Sprintf("BP %04d AB", rand.Intn(9999)),
			JenisKendaraan:     "Motor",
			StatusKetersediaan: "Tersedia",
		}

		if err := db.Create(&karyawan).Error; err != nil {
			log.Printf("Seeder: Gagal membuat data karyawan %s: %v", nama, err)
		}
	}

	log.Println("Seeder: Berhasil menabur data dummy Karyawan! 🛵")
}
