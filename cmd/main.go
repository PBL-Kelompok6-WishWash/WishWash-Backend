package main

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/controller"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/middleware"
	"net/http"
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
	// Initialize API versioning
	api := r.Group("/api/v1")
	{
		// Public routes: No authentication required
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected routes: Requires valid JWT token
		protected := api.Group("/protected")
		protected.Use(middleware.JWTAuthMiddleware()) 
		{
			// Test endpoint to verify JWT claims extraction
			protected.GET("/profil-saya", func(c *gin.Context) {
				userID, _ := c.Get("id_user")
				username, _ := c.Get("username")
				
				c.JSON(http.StatusOK, gin.H{
					"status":  "success",
					"message": "Autentikasi berhasil, akses ke rute terproteksi diizinkan",
					"data": gin.H{
						"user_id":  userID,
						"username": username,
					},
				})
			})
			
			// TODO: Add other protected routes (e.g., layananController, transaksiController)
		}
	}

	// 6. Buka restoran di port 8080
	log.Println("🚀 Server WishWash API berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}