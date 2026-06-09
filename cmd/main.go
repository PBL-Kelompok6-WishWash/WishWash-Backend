package main

import (
	"log"
	"strings"
	"time"
	// "net/http"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/controller"
	"github.com/PBL-Kelompok6-WishWash/backend/middleware"
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/route"
	"github.com/PBL-Kelompok6-WishWash/backend/seeder"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 1. Nyalakan Mesin Database
	config.ConnectDatabase()
	seeder.RunAllSeeders(config.DB)
	fixInconsistentOrders(config.DB)

	// 2. Pekerjakan "Koki" (Repository)
	userRepo := repository.NewUserRepository(config.DB)
	pelangganRepo := repository.NewPelangganRepository(config.DB)
	karyawanRepo := repository.NewKaryawanRepository(config.DB)
	adminRepo := repository.NewAdminRepository(config.DB)
	layananRepo := repository.NewLayananRepository(config.DB)
	parfumRepo := repository.NewParfumRepository(config.DB)
	promoRepo := repository.NewPromoRepository(config.DB)
	alamatRepo := repository.NewAlamatRepository(config.DB)
	notifikasiRepo := repository.NewNotifikasiRepository(config.DB)
	orderRepo := repository.NewOrderRepository(config.DB)
	chatRepo := repository.NewChatRepository(config.DB)

	// 3. Pekerjakan "Pelayan" (Controller)
	authController := controller.NewAuthController(userRepo, pelangganRepo, karyawanRepo, adminRepo, notifikasiRepo)
	pelangganController := controller.NewPelangganController(pelangganRepo, userRepo)
	karyawanController := controller.NewKaryawanController(karyawanRepo, userRepo)
	profileController := controller.NewProfileController(userRepo, adminRepo, karyawanRepo, pelangganRepo, alamatRepo)
	layananController := controller.NewLayananController(layananRepo)
	parfumController := controller.NewParfumController(parfumRepo)
	promoController := controller.NewPromoController(promoRepo)
	metodePembayaranController := controller.NewMetodePembayaranController(config.DB)
	alamatController := controller.NewAlamatController(alamatRepo, pelangganRepo)
	orderController := controller.NewOrderController(orderRepo, pelangganRepo, karyawanRepo, notifikasiRepo)
	chatController := controller.NewChatController(chatRepo, notifikasiRepo)
	penilaianController := controller.NewPenilaianController(config.DB)
	notifikasiController := controller.NewNotifikasiController(notifikasiRepo)

	// 4. Buka "Pintu Depan" menggunakan Gin Router
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Menyajikan folder uploads secara statis agar bisa diakses via URL/path
	r.Static("/uploads", "./uploads")

	// 5. Atur Papan Petunjuk Jalan (Routing API)
	api := r.Group("/api/v1")

	route.SetupAuthRoutes(api, authController)
	// route.SetupChatRoutes(api, chatController)

	profileRoutes := api.Group("/profile")
	profileRoutes.Use(middleware.JWTAuthMiddleware()) // Satpam 1 (Cek Token)
	{
		profileRoutes.GET("", profileController.GetProfile)
		profileRoutes.PUT("/update", profileController.UpdateProfile)
		profileRoutes.PUT("/password", profileController.UpdatePassword)
	}

	// Rute Alamat (Pelanggan Only)
	alamatRoutes := api.Group("/alamat")
	alamatRoutes.Use(middleware.JWTAuthMiddleware()) // Satpam 1 (Cek Token)
	{
		alamatRoutes.GET("", alamatController.GetAlamatPelanggan)
		alamatRoutes.POST("", alamatController.CreateAlamat)
		alamatRoutes.PUT("/:id", alamatController.UpdateAlamat)
		alamatRoutes.PUT("/:id/primary", alamatController.SetPrimaryAlamat)
		alamatRoutes.DELETE("/:id", alamatController.DeleteAlamat)
	}

	// Rute Order
	orderRoutes := api.Group("/order")
	orderRoutes.Use(middleware.JWTAuthMiddleware())
	{
		orderRoutes.GET("", orderController.GetOrdersPelanggan)
		orderRoutes.POST("", orderController.CreateOrder)
		orderRoutes.GET("/revenue", orderController.GetRevenueSummary)
		orderRoutes.GET("/by-kode/:kode", orderController.GetOrderByKode)
		orderRoutes.GET("/:id", orderController.GetOrderByID)
		orderRoutes.PUT("/:id", orderController.UpdateOrder)
		orderRoutes.POST("/scan-qr", orderController.ScanQR)
		orderRoutes.POST("/:id/penilaian", penilaianController.RateOrder)
	}

	// Alias untuk /orders sesuai dengan request
	ordersRoutes := api.Group("/orders")
	ordersRoutes.Use(middleware.JWTAuthMiddleware())
	{
		ordersRoutes.POST("/scan-qr", orderController.ScanQR)
	}

	// Rute Notifikasi untuk semua user terotentikasi (Pelanggan, Karyawan, Admin)
	notifikasiRoutes := api.Group("/notifikasi")
	{
		// WebSocket notifikasi tanpa middleware
		notifikasiRoutes.GET("/ws", notifikasiController.HandleNotifWS)

		// Rute terproteksi JWT
		notifikasiRoutes.Use(middleware.JWTAuthMiddleware())
		notifikasiRoutes.GET("", notifikasiController.GetNotifications)
		notifikasiRoutes.PUT("/:id/read", notifikasiController.MarkAsRead)
		notifikasiRoutes.PUT("/read-all", notifikasiController.MarkAllAsRead)
	}

	// Rute Layanan Pelanggan (General Authenticated Users)
	layananPubRoutes := api.Group("/layanan")
	layananPubRoutes.Use(middleware.JWTAuthMiddleware())
	{
		layananPubRoutes.GET("", layananController.GetAll)
	}

	// Rute Pelanggan Terotentikasi (General Authenticated Users, e.g. Karyawan/Kasir)
	pelangganPubRoutes := api.Group("/pelanggan")
	pelangganPubRoutes.Use(middleware.JWTAuthMiddleware())
	{
		pelangganPubRoutes.GET("", pelangganController.GetAll)
	}

	// Rute Promo Pelanggan (General Authenticated Users)
	promoPubRoutes := api.Group("/promo")
	promoPubRoutes.Use(middleware.JWTAuthMiddleware())
	{
		promoPubRoutes.GET("", promoController.GetAll)
	}

	// Rute Parfum Pelanggan (General Authenticated Users)
	parfumPubRoutes := api.Group("/parfum")
	parfumPubRoutes.Use(middleware.JWTAuthMiddleware())
	{
		parfumPubRoutes.GET("", parfumController.GetAll)
	}

	// chatRoutes := api.Group("/chat")
	// chatRoutes.Use(middleware.JWTAuthMiddleware()) // Dipasang satpam token biar aman
	// {
	// 	chatRoutes.GET("/room/:id_room_chat/messages", chatController.GetMessages)

	// 	chatRoutes.GET("/room/:id_room_chat/ws", chatController.HandleWS)
	// }

	// 🟢 REVISI YANG BENER (Ganti GetOrdersPelanggan jadi GetMessages):
    chatRoutes := api.Group("/chat")
    {
        // 1. Ambil history chat lama tetap dijaga satpam JWT, tapi pake fungsi asli kelompokmu (GetMessages)
        chatRoutes.GET("/room/:id_room_chat/messages", middleware.JWTAuthMiddleware(), chatController.GetMessages) 

        // 2. Jalur WebSocket dilepas dulu satpamnya biar jabat tangan dari HP lancar jaya
        chatRoutes.GET("/room/:id_room_chat/ws", chatController.HandleWS)

        // 3. Ambil daftar room chat aktif (dipakai di chat_screen.dart)
        chatRoutes.GET("/rooms", middleware.JWTAuthMiddleware(), chatController.GetRooms)

        // 4. Dapatkan atau buat room chat baru dari ID Order (dipakai di order_detail_screen.dart)
        chatRoutes.GET("/room/order/:id_order", middleware.JWTAuthMiddleware(), chatController.GetOrCreateRoom)
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

		// Rute Notifikasi
		adminRoutes.GET("/notifikasi", notifikasiController.GetNotifications)
		adminRoutes.PUT("/notifikasi/:id/read", notifikasiController.MarkAsRead)
		adminRoutes.PUT("/notifikasi/read-all", notifikasiController.MarkAllAsRead)
	}

	// 6. Buka restoran di port 8080
	log.Println("🚀 Server WishWash API berjalan di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}

func fixInconsistentOrders(db *gorm.DB) {
	log.Println("🔧 Memulai pemeriksaan konsistensi data pesanan di database...")
	var orders []model.Order
	err := db.Preload("RiwayatStatusDetail.ReferensiStatus").Find(&orders).Error
	if err != nil {
		log.Printf("❌ Gagal mengambil data pesanan untuk koreksi: %v", err)
		return
	}

	correctedCount := 0
	for _, order := range orders {
		if order.Kuantitas <= 0 {
			continue
		}

		// Cari status terakhir
		currentStatus := "Pesanan Diterima"
		maxUrutan := 0
		for _, rs := range order.RiwayatStatusDetail {
			if rs.ReferensiStatus.UrutanTahap > maxUrutan {
				maxUrutan = rs.ReferensiStatus.UrutanTahap
				currentStatus = rs.ReferensiStatus.NamaStatus
			}
		}

		// Jika status saat ini masih di bawah/sama dengan Proses Timbang
		lowerCurrent := strings.ToLower(currentStatus)
		if strings.Contains(lowerCurrent, "timbang") || strings.Contains(lowerCurrent, "jemput") || strings.Contains(lowerCurrent, "terima") {
			// Cari status Proses Timbang untuk mendapatkan ID-nya jika diperlukan
			var refStatuses []model.ReferensiStatusLayanan
			db.Where("id_layanan = ?", order.LayananID).Order("urutan_tahap asc").Find(&refStatuses)

			timbangIdx := -1
			for i, ref := range refStatuses {
				refNameLower := strings.ToLower(ref.NamaStatus)
				if strings.Contains(refNameLower, "timbang") {
					timbangIdx = i
					break
				}
			}

			// Jika Proses Timbang ditemukan, cari status setelahnya
			if timbangIdx != -1 && timbangIdx < len(refStatuses)-1 {
				nextStatus := refStatuses[timbangIdx+1]
				
				// Pastikan status tersebut belum ada di riwayat agar tidak duplikat
				alreadyHasNext := false
				for _, rs := range order.RiwayatStatusDetail {
					if rs.ReferensiStatusID == nextStatus.IDReferensiStatus {
						alreadyHasNext = true
						break
					}
				}

				if !alreadyHasNext {
					newHistory := model.RiwayatStatusDetail{
						ReferensiStatusID: nextStatus.IDReferensiStatus,
						OrderID:           order.IDOrder,
						KaryawanID:        nil,
						WaktuUpdate:       time.Now(),
					}
					if err := db.Create(&newHistory).Error; err == nil {
						correctedCount++
						log.Printf("✅ Koreksi Pesanan #%s: Status dimajukan dari '%s' ke '%s' karena berat sudah diisi (%v kg)", 
							order.KodeOrder, currentStatus, nextStatus.NamaStatus, order.Kuantitas)
					}
				}
			}
		}
	}
	log.Printf("🔧 Pemeriksaan selesai. Berhasil mengoreksi %d data pesanan.", correctedCount)
}
