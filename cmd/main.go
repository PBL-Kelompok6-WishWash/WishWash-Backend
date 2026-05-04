package main

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/controller"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/seeder"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Nyalakan Mesin Database (Sudah disesuaikan dengan nama fungsimu!)
	config.ConnectDatabase()

	seeder.RunAllSeeders(config.DB)
	// 2. Pekerjakan "Koki" (Repository) dan beri dia akses ke Database
	// ⚠️ Pastikan di dalam file config kamu benar-benar ada variabel global bernama 'DB'
	userRepo := repository.NewUserRepository(config.DB)

	// 3. Pekerjakan "Pelayan" (Controller) dan sambungkan dia dengan sang Koki
	authController := controller.NewAuthController(userRepo)

	// 4. Buka "Pintu Depan" menggunakan Gin Router
	r := gin.Default()

	// 5. Atur Papan Petunjuk Jalan (Routing API)
	api := r.Group("/api/v1")
	{
		// Modul Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
	}

	// 6. Buka restoran di port 8080
	log.Println("🚀 Server WishWash API berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}