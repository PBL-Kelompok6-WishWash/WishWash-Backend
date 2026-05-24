package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
}

type orderController struct {
	orderRepo     repository.OrderRepository
	pelangganRepo repository.PelangganRepository
}

func NewOrderController(oRepo repository.OrderRepository, pRepo repository.PelangganRepository) OrderController {
	return &orderController{oRepo, pRepo}
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
	pelangganID, err := ctrl.getPelangganIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	orders, err := ctrl.orderRepo.FindAllByPelangganID(pelangganID)
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
		AlamatPengambilanID uint    `json:"id_alamat_pengambilan" binding:"required"`
		AlamatPenyerahanID  *uint   `json:"id_alamat_penyerahan"`
		ParfumID            uint    `json:"id_parfum" binding:"required"`
		LayananID           uint    `json:"id_layanan" binding:"required"`
		KeteranganLokasi    string  `json:"keterangan_lokasi"`
		JadwalPickup        string  `json:"jadwal_pickup" binding:"required"` // Format: YYYY-MM-DD HH:MM
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

	if roleID == 1 || roleID == 2 {
		if input.PelangganID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID Pelanggan wajib diisi untuk Karyawan/Admin"})
			return
		}
		pelangganID = *input.PelangganID
	} else {
		pelangganID, err = ctrl.getPelangganIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
			return
		}
	}

	// Parse jadwal pickup
	pickupTime, err := time.Parse("2006-01-02 15:04", input.JadwalPickup)
	if err != nil {
		// Fallback to try RFC3339
		pickupTime, err = time.Parse(time.RFC3339, input.JadwalPickup)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format jadwal pickup salah. Harus YYYY-MM-DD HH:MM"})
			return
		}
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
