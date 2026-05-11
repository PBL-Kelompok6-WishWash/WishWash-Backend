package main

import (
	"log"
	"net/http"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/controller"
	"github.com/PBL-Kelompok6-WishWash/backend/middleware"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/seeder"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Nyalakan Mesin Database
	config.ConnectDatabase()
	seeder.RunAllSeeders(config.DB)

	// 2. Pekerjakan "Koki" (Repository)
	userRepo := repository.NewUserRepository(config.DB)
	pelangganRepo := repository.NewPelangganRepository(config.DB)
	karyawanRepo := repository.NewKaryawanRepository(config.DB)
	adminRepo := repository.NewAdminRepository(config.DB)

	// 3. Pekerjakan "Pelayan" (Controller)
	authController := controller.NewAuthController(userRepo, pelangganRepo, karyawanRepo, adminRepo)
	pelangganController := controller.NewPelangganController(pelangganRepo, userRepo)
	// 💡 TAMBAHAN BARU: Inisialisasi Profile Controller (Satu pelayan untuk semua role)
	profileController := controller.NewProfileController(userRepo, adminRepo, karyawanRepo, pelangganRepo)

	// 4. Buka "Pintu Depan" menggunakan Gin Router
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 5. Atur Papan Petunjuk Jalan (Routing API)
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
			// Test endpoint
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
			
			// Rute untuk Update Profile (Pakai PUT karena meng-update data yang sudah ada)
			protected.PUT("/profile/update", profileController.UpdateProfile)
			protected.PUT("/password/update", profileController.UpdatePassword)
			
			protected.GET("/pelanggan", pelangganController.GetAll)
			protected.GET("/pelanggan/:id", pelangganController.GetByID)
			protected.POST("/pelanggan", pelangganController.Create)
			protected.PUT("/pelanggan/:id", pelangganController.Update)
			protected.DELETE("/pelanggan/:id", pelangganController.Delete)
		}
	}

	// 6. Buka restoran di port 8080
	log.Println("🚀 Server WishWash API berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}