package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
)

func generateKodeOrder() string {
	rand.Seed(time.Now().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return fmt.Sprintf("WW-%s", string(result))
}

type OrderController interface {
	GetOrdersPelanggan(c *gin.Context)
	CreateOrder(c *gin.Context)
	GetOrderByID(c *gin.Context)
	GetOrderByKode(c *gin.Context)
	UpdateOrder(c *gin.Context)
	ScanQR(c *gin.Context)
	GetRevenueSummary(c *gin.Context)
}

type orderController struct {
	orderRepo     repository.OrderRepository
	pelangganRepo repository.PelangganRepository
	karyawanRepo  repository.KaryawanRepository
}

func NewOrderController(oRepo repository.OrderRepository, pRepo repository.PelangganRepository, kRepo repository.KaryawanRepository) OrderController {
	return &orderController{oRepo, pRepo, kRepo}
}

func (ctrl *orderController) getPelangganIDFromContext(c *gin.Context) (uint, error) {
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		return 0, http.ErrNoCookie
	}
	userID := uint(userIDFloat.(float64))
	pelanggan, err := ctrl.pelangganRepo.FindByUserID(userID)
	if err != nil {
		return 0, err
	}
	return pelanggan.IDPelanggan, nil
}

func (ctrl *orderController) GetOrdersPelanggan(c *gin.Context) {
	roleData, exists := c.Get("id_role")
	roleID := 3 // default to customer
	if exists {
		roleID = int(roleData.(float64))
	}

	var orders []model.Order
	var err error

	if roleID == 1 || roleID == 2 {
		// Admin atau Karyawan: ambil semua order
		orders, err = ctrl.orderRepo.FindAll()
	} else {
		// Pelanggan: ambil order milik pelanggan itu saja
		pelangganID, errPelanggan := ctrl.getPelangganIDFromContext(c)
		if errPelanggan != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
			return
		}
		orders, err = ctrl.orderRepo.FindAllByPelangganID(pelangganID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pesanan berhasil diambil",
		"data":    orders,
	})
}

func (ctrl *orderController) CreateOrder(c *gin.Context) {
	var pelangganID uint
	var err error

	roleData, exists := c.Get("id_role")
	roleID := 3 // default to customer
	if exists {
		roleID = int(roleData.(float64))
	}

	var input struct {
		PelangganID         *uint   `json:"id_pelanggan"`
		PaketLayananID      *uint   `json:"id_paket_layanan"`
		AlamatPengambilanID *uint   `json:"id_alamat_pengambilan"`
		AlamatPenyerahanID  *uint   `json:"id_alamat_penyerahan"`
		ParfumID            uint    `json:"id_parfum" binding:"required"`
		LayananID           uint    `json:"id_layanan" binding:"required"`
		KeteranganLokasi    string  `json:"keterangan_lokasi"`
		JadwalPickup        string  `json:"jadwal_pickup"` // Format: YYYY-MM-DD HH:MM
		TipeLogistik        string  `json:"tipe_logistik" binding:"required"`
		HargaSaatIni        float64 `json:"harga_saat_ini" binding:"required"`
		Kuantitas           float64 `json:"kuantitas"`
		TotalBayar          float64 `json:"total_bayar"`
		CatatanOrder        string  `json:"catatan_order"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	if input.TipeLogistik == "Courier Delivery" {
		if input.AlamatPengambilanID == nil || *input.AlamatPengambilanID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat pengambilan wajib diisi untuk pengiriman kurir"})
			return
		}
		if input.JadwalPickup == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Jadwal pickup wajib diisi untuk pengiriman kurir"})
			return
		}
	}

	var karyawanID *uint
	if roleID == 1 || roleID == 2 {
		if input.PelangganID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID Pelanggan wajib diisi untuk Karyawan/Admin"})
			return
		}
		pelangganID = *input.PelangganID

		userIDFloat, exists := c.Get("id_user")
		if exists {
			userID := uint(userIDFloat.(float64))
			karyawan, errKaryawan := ctrl.karyawanRepo.FindByUserID(userID)
			if errKaryawan == nil {
				karyawanID = &karyawan.IDKaryawan
			}
		}
	} else {
		pelangganID, err = ctrl.getPelangganIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
			return
		}
	}

	var pickupTime *time.Time
	if input.JadwalPickup != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04", input.JadwalPickup)
		if err != nil {
			// Fallback to try RFC3339
			parsedTime, err = time.Parse(time.RFC3339, input.JadwalPickup)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format jadwal pickup salah. Harus YYYY-MM-DD HH:MM"})
				return
			}
		}
		pickupTime = &parsedTime
	}

	order := model.Order{
		PelangganID:         pelangganID,
		KodeOrder:           generateKodeOrder(),
		PaketLayananID:      input.PaketLayananID,
		AlamatPengambilanID: input.AlamatPengambilanID,
		AlamatPenyerahanID:  input.AlamatPenyerahanID,
		ParfumID:            input.ParfumID,
		LayananID:           input.LayananID,
		KeteranganLokasi:    input.KeteranganLokasi,
		JadwalPickup:        pickupTime,
		TipeLogistik:        input.TipeLogistik,
		HargaSaatIni:        input.HargaSaatIni,
		Kuantitas:           input.Kuantitas,
		TotalBayar:          input.TotalBayar,
		CatatanOrder:        input.CatatanOrder,
		TglPesanan:          time.Now(),
		KaryawanID:          karyawanID,
	}

	if err := ctrl.orderRepo.Create(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesanan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesanan berhasil dibuat",
		"data":    order,
	})
}

func (ctrl *orderController) GetOrderByID(c *gin.Context) {
	idOrderParam := c.Param("id")
	idOrder, err := strconv.ParseUint(idOrderParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Order tidak valid"})
		return
	}

	order, err := ctrl.orderRepo.FindByID(uint(idOrder))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pesanan berhasil diambil",
		"data":    order,
	})
}

func (ctrl *orderController) GetOrderByKode(c *gin.Context) {
	kodeOrder := c.Param("kode")
	if kodeOrder == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode order tidak valid"})
		return
	}

	order, err := ctrl.orderRepo.FindByKodeOrder(kodeOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pesanan berhasil diambil",
		"data":    order,
	})
}

func (ctrl *orderController) UpdateOrder(c *gin.Context) {
	idOrderParam := c.Param("id")
	idOrder, err := strconv.ParseUint(idOrderParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Order tidak valid"})
		return
	}

	// 1. Dapatkan KaryawanID dari UserID yang terotentikasi (jika pengguna adalah Karyawan)
	var karyawanID *uint
	userIDFloat, exists := c.Get("id_user")
	if exists {
		userID := uint(userIDFloat.(float64))
		karyawan, errKaryawan := ctrl.karyawanRepo.FindByUserID(userID)
		if errKaryawan == nil {
			karyawanID = &karyawan.IDKaryawan
		}
	}

	// 2. Ambil data order existing
	order, err := ctrl.orderRepo.FindByID(uint(idOrder))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	// 3. Bind input JSON
	var input struct {
		Status              string   `json:"status"`            // Contoh: "Diproses", "Selesai"
		Kuantitas           *float64 `json:"kuantitas"`         // Contoh: 3.5 (dalam kg)
		TotalBayar          *float64 `json:"total_bayar"`       // Total bayar baru jika diubah
		StatusPembayaran    string   `json:"status_pembayaran"` // Contoh: "Paid", "Lunas", "Unpaid"
		MetodeBayar         string   `json:"metode_bayar"`      // Contoh: "Cash", "QRIS"
		TipeLogistik        string   `json:"tipe_logistik"`
		AlamatPenyerahanID  *uint    `json:"id_alamat_penyerahan"`
		AlamatPengambilanID *uint    `json:"id_alamat_pengambilan"`
		CatatanOrder        string   `json:"catatan_order"`
		IsCourierOnWay      *bool    `json:"is_courier_on_way"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// 4. Update status jika dikirim
	if input.Status != "" {
		var refStatus model.ReferensiStatusLayanan
		var errRef error
		
		if input.Status == "Dibatalkan" || input.Status == "Batal" {
			errRef = config.DB.Where("id_layanan = ? AND nama_status = ?", order.LayananID, "Dibatalkan").First(&refStatus).Error
			if errRef != nil {
				var maxUrutan int
				config.DB.Model(&model.ReferensiStatusLayanan{}).Where("id_layanan = ?", order.LayananID).Select("COALESCE(MAX(urutan_tahap), 0)").Row().Scan(&maxUrutan)
				refStatus = model.ReferensiStatusLayanan{
					LayananID:   order.LayananID,
					NamaStatus:  "Dibatalkan",
					UrutanTahap: maxUrutan + 1,
				}
				config.DB.Create(&refStatus)
				errRef = nil
			}
		} else {
			errRef = config.DB.Where("id_layanan = ? AND nama_status = ?", order.LayananID, input.Status).First(&refStatus).Error
			if errRef != nil {
				errRef = config.DB.Where("id_layanan = ? AND LOWER(nama_status) = LOWER(?)", order.LayananID, input.Status).First(&refStatus).Error
			}
		}

		if errRef == nil {
			// Buat riwayat status detail baru
			history := model.RiwayatStatusDetail{
				ReferensiStatusID: refStatus.IDReferensiStatus,
				OrderID:           order.IDOrder,
				KaryawanID:        karyawanID,
				WaktuUpdate:       time.Now(),
			}
			if err := config.DB.Create(&history).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan riwayat status: " + err.Error()})
				return
			}

			// If the order is transitioning away from "Pesanan Diterima" (i.e. it was accepted by the employee)
			isCurrentDiterima := false
			if len(order.RiwayatStatusDetail) > 0 {
				var maxUrutan int = 0
				var currentStatusName string
				for _, h := range order.RiwayatStatusDetail {
					if h.ReferensiStatus.UrutanTahap > maxUrutan {
						maxUrutan = h.ReferensiStatus.UrutanTahap
						currentStatusName = h.ReferensiStatus.NamaStatus
					}
				}
				if currentStatusName == "Pesanan Diterima" || maxUrutan == 1 {
					isCurrentDiterima = true
				}
			} else {
				isCurrentDiterima = true
			}

			if isCurrentDiterima {
				var diterimaStatus model.ReferensiStatusLayanan
				errDiterima := config.DB.Where("id_layanan = ? AND urutan_tahap = 1", order.LayananID).First(&diterimaStatus).Error
				if errDiterima == nil {
					config.DB.Model(&model.RiwayatStatusDetail{}).
						Where("order_id = ? AND referensi_status_id = ?", order.IDOrder, diterimaStatus.IDReferensiStatus).
						Update("waktu_update", time.Now())
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status '" + input.Status + "' tidak ditemukan untuk layanan ini"})
			return
		}
	}

	// 5. Update kuantitas/weight jika dikirim
	if input.Kuantitas != nil {
		order.Kuantitas = *input.Kuantitas
		// Hitung ulang total bayar jika input.TotalBayar tidak dikirim secara manual
		if input.TotalBayar != nil {
			order.TotalBayar = *input.TotalBayar
		} else {
			order.TotalBayar = order.Kuantitas * order.HargaSaatIni
		}
	} else if input.TotalBayar != nil {
		order.TotalBayar = *input.TotalBayar
	}

	if input.TipeLogistik != "" {
		order.TipeLogistik = input.TipeLogistik
		if input.TipeLogistik == "Drop-off" {
			order.AlamatPenyerahanID = nil
		}
	}
	if input.AlamatPenyerahanID != nil {
		order.AlamatPenyerahanID = input.AlamatPenyerahanID
	}
	if input.AlamatPengambilanID != nil {
		order.AlamatPengambilanID = input.AlamatPengambilanID
	}
	if input.CatatanOrder != "" {
		order.CatatanOrder = input.CatatanOrder
	}
	if input.IsCourierOnWay != nil {
		order.IsCourierOnWay = *input.IsCourierOnWay
		if *input.IsCourierOnWay {
			var latestHistory model.RiwayatStatusDetail
			errLatest := config.DB.Preload("ReferensiStatus").
				Where("order_id = ?", order.IDOrder).
				Order("id_riwayat_status_detail desc").
				First(&latestHistory).Error
			
			if errLatest == nil && latestHistory.ReferensiStatus.IDReferensiStatus != 0 {
				statusName := latestHistory.ReferensiStatus.NamaStatus
				if statusName == "Penjemputan" || statusName == "Siap Diantar" {
					config.DB.Model(&latestHistory).Update("waktu_update", time.Now())
				}
			}
		}
	}

	// Jika Karyawan meng-update order, set KaryawanID di order
	if karyawanID != nil {
		order.KaryawanID = karyawanID
	}

	// Simpan perubahan order ke DB
	if err := ctrl.orderRepo.Update(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate order: " + err.Error()})
		return
	}

	// 6. Update Pembayaran jika statusPembayaran dikirim
	if input.StatusPembayaran != "" {
		var pembayaran model.Pembayaran
		errPem := config.DB.Where("id_order = ?", order.IDOrder).First(&pembayaran).Error
		if errPem == nil {
			// Update pembayaran yang ada
			pembayaran.StatusPembayaran = input.StatusPembayaran
			pembayaran.KaryawanID = karyawanID
			if input.MetodeBayar != "" {
				pembayaran.MetodeBayar = input.MetodeBayar
			}
			if input.TotalBayar != nil {
				pembayaran.JumlahBayar = *input.TotalBayar
			} else {
				pembayaran.JumlahBayar = order.TotalBayar
			}
			pembayaran.TglPembayaran = time.Now()
			config.DB.Save(&pembayaran)
		} else {
			// Buat pembayaran baru
			metode := "Cash"
			if input.MetodeBayar != "" {
				metode = input.MetodeBayar
			}
			jumlah := order.TotalBayar
			if input.TotalBayar != nil {
				jumlah = *input.TotalBayar
			}
			pembayaran = model.Pembayaran{
				OrderID:          order.IDOrder,
				KaryawanID:       karyawanID,
				MetodeBayar:      metode,
				JumlahBayar:      jumlah,
				StatusPembayaran: input.StatusPembayaran,
				TglPembayaran:    time.Now(),
			}
			config.DB.Create(&pembayaran)
		}
	}

	// Fetch kembali data lengkap untuk dikembalikan ke client
	updatedOrder, _ := ctrl.orderRepo.FindByID(order.IDOrder)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesanan berhasil diperbarui",
		"data":    updatedOrder,
	})
}

func (ctrl *orderController) ScanQR(c *gin.Context) {
	var input struct {
		OrderID interface{} `json:"order_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Input tidak valid: " + err.Error(), "error": err.Error()})
		return
	}

	var idOrder uint
	switch v := input.OrderID.(type) {
	case float64:
		idOrder = uint(v)
	case string:
		parsed, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format order_id tidak valid", "error": err.Error()})
			return
		}
		idOrder = uint(parsed)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format order_id tidak valid", "error": "Invalid type"})
		return
	}

	// 1. Ambil data order
	order, err := ctrl.orderRepo.FindByID(idOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Pesanan tidak ditemukan", "error": err.Error()})
		return
	}

	// 2. Dapatkan KaryawanID (jika ada)
	var karyawanID *uint
	userIDFloat, exists := c.Get("id_user")
	if exists {
		userID := uint(userIDFloat.(float64))
		karyawan, errKaryawan := ctrl.karyawanRepo.FindByUserID(userID)
		if errKaryawan == nil {
			karyawanID = &karyawan.IDKaryawan
		}
	}

	// 3. Cari status terakhir dan periksa status akhir
	var currentUrutan int = 0
	var lastStatusName string
	if len(order.RiwayatStatusDetail) > 0 {
		maxUrutanFound := 0
		for _, rs := range order.RiwayatStatusDetail {
			if rs.ReferensiStatus.UrutanTahap > currentUrutan {
				currentUrutan = rs.ReferensiStatus.UrutanTahap
			}
			if rs.ReferensiStatus.UrutanTahap > maxUrutanFound {
				maxUrutanFound = rs.ReferensiStatus.UrutanTahap
				lastStatusName = rs.ReferensiStatus.NamaStatus
			}
		}
	}
	// Jika status terakhir merupakan status final, blok update
	finalStatuses := map[string]bool{"Selesai": true, "Batal": true, "Dibatalkan": true}
	if finalStatuses[lastStatusName] {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Pesanan sudah selesai", "error": "order completed"})
		return
	}

	// 4. Pastikan tidak di status akhir (fallback safety)
	var maxUrutan int
	config.DB.Model(&model.ReferensiStatusLayanan{}).Where("id_layanan = ?", order.LayananID).Select("COALESCE(MAX(urutan_tahap),0)").Row().Scan(&maxUrutan)
	if currentUrutan >= maxUrutan {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "Pesanan sudah selesai", "error": "order completed"})
		return
	}
	var nextStatus model.ReferensiStatusLayanan
	err = config.DB.Where("id_layanan = ? AND urutan_tahap = ?", order.LayananID, currentUrutan+1).First(&nextStatus).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Pesanan sudah berada di status akhir atau tidak dapat dilanjutkan", "error": err.Error()})
		return
	}

	// 5. Update status
	history := model.RiwayatStatusDetail{
		ReferensiStatusID: nextStatus.IDReferensiStatus,
		OrderID:           order.IDOrder,
		KaryawanID:        karyawanID,
		WaktuUpdate:       time.Now(),
	}
	if err := config.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengupdate status pesanan", "error": err.Error()})
		return
	}

	updatedOrder, _ := ctrl.orderRepo.FindByID(order.IDOrder)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status pesanan berhasil diperbarui menjadi " + nextStatus.NamaStatus,
		"data":    updatedOrder,
	})
}

func (ctrl *orderController) GetRevenueSummary(c *gin.Context) {
	roleData, exists := c.Get("id_role")
	roleID := 3 // default to customer
	if exists {
		roleID = int(roleData.(float64))
	}

	if roleID != 1 && roleID != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Admin atau Karyawan yang dapat mengakses data pendapatan"})
		return
	}

	var payments []model.Pembayaran
	err := config.DB.
		Preload("Order.Pelanggan").
		Preload("Order.Layanan").
		Where("LOWER(status_pembayaran) IN (?, ?)", "paid", "lunas").
		Order("tgl_pembayaran desc").
		Find(&payments).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pendapatan: " + err.Error()})
		return
	}

	now := time.Now()
	loc := now.Location()

	targetMonth := int(now.Month())
	targetYear := now.Year()

	if mStr := c.Query("month"); mStr != "" {
		if mVal, err := strconv.Atoi(mStr); err == nil && mVal >= 1 && mVal <= 12 {
			targetMonth = mVal
		}
	}
	if yStr := c.Query("year"); yStr != "" {
		if yVal, err := strconv.Atoi(yStr); err == nil {
			targetYear = yVal
		}
	}

	year, month, day := now.Date()
	todayStart := time.Date(year, month, day, 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	yesterdayStart := todayStart.Add(-24 * time.Hour)
	yesterdayEnd := todayStart

	monthStart := time.Date(targetYear, time.Month(targetMonth), 1, 0, 0, 0, 0, loc)
	monthEnd := monthStart.AddDate(0, 1, 0)

	var todayRevenue float64
	var yesterdayRevenue float64
	var monthlyRevenue float64
	var accumulatedRevenue float64

	type TransactionItem struct {
		ID            string    `json:"id"`
		MethodType    string    `json:"method_type"`
		Title         string    `json:"title"`
		Subtitle      string    `json:"subtitle"`
		Time          time.Time `json:"time"`
		Amount        float64   `json:"amount"`
		PaymentMethod string    `json:"payment_method"`
		FotoPelanggan string    `json:"foto_pelanggan"`
	}

	var transactions []TransactionItem

	for _, p := range payments {
		accumulatedRevenue += p.JumlahBayar

		pTimeLocal := p.TglPembayaran.In(loc)

		// Check if today
		if pTimeLocal.After(todayStart) && pTimeLocal.Before(todayEnd) {
			todayRevenue += p.JumlahBayar
		}

		// Check if yesterday
		if pTimeLocal.After(yesterdayStart) && pTimeLocal.Before(yesterdayEnd) {
			yesterdayRevenue += p.JumlahBayar
		}

		// Check if this month
		if pTimeLocal.After(monthStart) && pTimeLocal.Before(monthEnd) {
			monthlyRevenue += p.JumlahBayar

			methodType := "digital"
			if p.MetodeBayar == "Tunai" || p.MetodeBayar == "Cash" {
				methodType = "cash"
			}

			title := fmt.Sprintf("Order #%s", p.Order.KodeOrder)
			if p.Order.KodeOrder == "" {
				title = fmt.Sprintf("Order #WW-%d", p.OrderID)
			}

			subtitle := fmt.Sprintf("%s • %s", p.Order.Pelanggan.NamaLengkap, p.Order.Layanan.NamaLayanan)

			transactions = append(transactions, TransactionItem{
				ID:            fmt.Sprintf("TX-%d", p.IDPembayaran),
				MethodType:    methodType,
				Title:         title,
				Subtitle:      subtitle,
				Time:          p.TglPembayaran,
				Amount:        p.JumlahBayar,
				PaymentMethod: p.MetodeBayar,
				FotoPelanggan: p.Order.Pelanggan.FotoPelanggan,
			})
		}
	}

	// Hitung persentase trend
	var trendText string
	if yesterdayRevenue == 0 {
		if todayRevenue > 0 {
			trendText = "+100%"
		} else {
			trendText = "0%"
		}
	} else {
		diff := ((todayRevenue - yesterdayRevenue) / yesterdayRevenue) * 100
		if diff >= 0 {
			trendText = fmt.Sprintf("+%.1f%%", diff)
		} else {
			trendText = fmt.Sprintf("%.1f%%", diff)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"today_revenue":       todayRevenue,
		"yesterday_revenue":   yesterdayRevenue,
		"monthly_revenue":     monthlyRevenue,
		"accumulated_revenue": accumulatedRevenue,
		"percentage_trend":    trendText,
		"transactions":        transactions,
	})
}

