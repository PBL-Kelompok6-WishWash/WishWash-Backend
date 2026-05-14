package main

import (
	"log"
	// "net/http"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/controller"
	"github.com/PBL-Kelompok6-WishWash/backend/middleware"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/route"
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
	layananRepo := repository.NewLayananRepository(config.DB)
	parfumRepo := repository.NewParfumRepository(config.DB)
	promoRepo := repository.NewPromoRepository(config.DB)

	// 3. Pekerjakan "Pelayan" (Controller)
	authController := controller.NewAuthController(userRepo, pelangganRepo, karyawanRepo, adminRepo)
	pelangganController := controller.NewPelangganController(pelangganRepo, userRepo)
	karyawanController := controller.NewKaryawanController(karyawanRepo, userRepo)
	profileController := controller.NewProfileController(userRepo, adminRepo, karyawanRepo, pelangganRepo)
	layananController := controller.NewLayananController(layananRepo)
	parfumController := controller.NewParfumController(parfumRepo)
	promoController := controller.NewPromoController(promoRepo)
	metodePembayaranController := controller.NewMetodePembayaranController(config.DB)

	// 4. Buka "Pintu Depan" menggunakan Gin Router
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// 5. Atur Papan Petunjuk Jalan (Routing API)
	api := r.Group("/api/v1")

	route.SetupAuthRoutes(api, authController)

	profileRoutes := api.Group("/profile")
    profileRoutes.Use(middleware.JWTAuthMiddleware()) // Satpam 1 (Cek Token)
    {
        profileRoutes.GET("", profileController.GetProfile)
        profileRoutes.PUT("/update", profileController.UpdateProfile)
        profileRoutes.PUT("/password", profileController.UpdatePassword)
    }

    // B. Rute Khusus Admin (Hanya Role 1 yang bisa akses)
    adminRoutes := api.Group("/admin")
    adminRoutes.Use(middleware.JWTAuthMiddleware(), middleware.AdminOnly()) // Satpam 1 & 2
    {
        adminRoutes.GET("/pelanggan", pelangganController.GetAll)
        adminRoutes.GET("/pelanggan/:id", pelangganController.GetByID)
        adminRoutes.POST("/pelanggan", pelangganController.Create)
        adminRoutes.PUT("/pelanggan/:id", pelangganController.Update)
        adminRoutes.DELETE("/pelanggan/:id", pelangganController.Delete)
        
        adminRoutes.GET("/karyawan", karyawanController.GetAll)
        adminRoutes.GET("/karyawan/:id", karyawanController.GetByID)
        adminRoutes.POST("/karyawan", karyawanController.Create)
        adminRoutes.PUT("/karyawan/:id", karyawanController.Update)
        adminRoutes.DELETE("/karyawan/:id", karyawanController.Delete)
        
        // Rute Layanan
        adminRoutes.GET("/layanan", layananController.GetAll)
        adminRoutes.GET("/layanan/:id", layananController.GetByID)
        adminRoutes.POST("/layanan", layananController.Create)
        adminRoutes.PUT("/layanan/:id", layananController.Update)
        adminRoutes.DELETE("/layanan/:id", layananController.Delete)
        
        // Rute Parfum
        adminRoutes.GET("/parfum", parfumController.GetAll)
        adminRoutes.GET("/parfum/:id", parfumController.GetByID)
        adminRoutes.POST("/parfum", parfumController.Create)
        adminRoutes.PUT("/parfum/:id", parfumController.Update)
        adminRoutes.DELETE("/parfum/:id", parfumController.Delete)

        // Rute Promo
        adminRoutes.GET("/promo", promoController.GetAll)
        adminRoutes.GET("/promo/:id", promoController.GetByID)
        adminRoutes.POST("/promo", promoController.Create)
        adminRoutes.PUT("/promo/:id", promoController.Update)
        adminRoutes.DELETE("/promo/:id", promoController.Delete)

        // Rute Metode Pembayaran
        adminRoutes.GET("/metode-pembayaran", metodePembayaranController.GetAll)
        adminRoutes.GET("/metode-pembayaran/:id", metodePembayaranController.GetByID)
        adminRoutes.POST("/metode-pembayaran", metodePembayaranController.Create)
        adminRoutes.PUT("/metode-pembayaran/:id", metodePembayaranController.Update)
        adminRoutes.DELETE("/metode-pembayaran/:id", metodePembayaranController.Delete)
    }

	// 6. Buka restoran di port 8080
	log.Println("🚀 Server WishWash API berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}