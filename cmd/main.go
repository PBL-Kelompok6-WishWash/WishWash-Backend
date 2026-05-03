package main

import (
	"fmt"
	"github.com/PBL-Kelompok6-WishWash/backend/config" // Sesuaikan nama modul
)

func main() {
	fmt.Println("🚀 Memulai server WishWash...")

	// Memanggil fungsi ConnectDatabase yang ada di folder config
	config.ConnectDatabase()
}