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

func SeedPelanggan(db *gorm.DB) {
	// 1. Cek apakah tabel pelanggan sudah ada isinya biar nggak duplikat pas di-run berkali-kali
	var count int64
	if err := db.Model(&model.Pelanggan{}).Count(&count).Error; err != nil {
		log.Printf("Gagal mengecek data pelanggan: %v", err)
		return
	}

	if count > 0 {
		log.Println("Seeder: Data Pelanggan sudah ada, skip proses seeding.")
		return
	}

	// 2. Siapkan password default (misal: "password123")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// 3. Buat Pelanggan 1 khusus
	hashedPasswordSpecific, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	user1 := model.User{
		Username: "pelanggan1",
		Email:    "pelanggan1@gmail.com",
		Password: string(hashedPasswordSpecific),
		RoleID:   3,
	}
	db.Create(&user1)
	pelanggan1 := model.Pelanggan{
		UserID:      user1.IDUser,
		NamaLengkap: "Pelanggan Spesifik 1",
		NoTelp:      "081111111111",
	}
	db.Create(&pelanggan1)

	firstNames := []string{"Budi", "Andi", "Roni", "Siti", "Ayu", "Joko", "Tono", "Dewi", "Rina", "Agus", "Bagus", "Cahyo", "Dina", "Eka", "Fajar", "Gita", "Hadi", "Indra", "Jamil", "Kartika"}
	lastNames := []string{"Santoso", "Wijaya", "Kurniawan", "Setiawan", "Pratama", "Putra", "Putri", "Sari", "Lestari", "Hidayat", "Nugroho", "Saputra", "Wahyudi", "Ramadhan"}

	rand.Seed(time.Now().UnixNano())

	// 4. Looping untuk insert ke Database 20 Pelanggan Acak
	for i := 1; i <= 20; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		nama := fmt.Sprintf("%s %s", firstName, lastName)
		username := strings.ToLower(fmt.Sprintf("%s%s%d", firstName, lastName, rand.Intn(100)))
		email := fmt.Sprintf("%s@gmail.com", username)
		noTelp := fmt.Sprintf("081234%06d", rand.Intn(999999))

		user := model.User{
			Username: username,
			Email:    email,
			Password: string(hashedPassword),
			RoleID:   3,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Seeder: Gagal membuat user untuk %s: %v", username, err)
			continue
		}

		// Insert Pelanggan pakai ID User yang baru dibuat
		pelanggan := model.Pelanggan{
			UserID:      user.IDUser,
			NamaLengkap: nama,
			NoTelp:      noTelp,
		}

		if err := db.Create(&pelanggan).Error; err != nil {
			log.Printf("Seeder: Gagal membuat data pelanggan %s: %v", nama, err)
		}
	}

	log.Println("Seeder: Berhasil menabur data dummy Pelanggan! 🌱")
}