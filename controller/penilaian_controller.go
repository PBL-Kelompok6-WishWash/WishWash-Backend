package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenilaianController struct {
	db *gorm.DB
}

func NewPenilaianController(db *gorm.DB) *PenilaianController {
	return &PenilaianController{db}
}

func (c *PenilaianController) RateOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	orderID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID order tidak valid"})
		return
	}

	// Verify order exists
	var order model.Order
	if err := c.db.Preload("Pelanggan").First(&order, orderID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order tidak ditemukan"})
		return
	}

	// Check if already rated
	var existing model.Penilaian
	if err := c.db.Where("id_order = ?", orderID).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Order ini sudah diberikan ulasan"})
		return
	}

	var input struct {
		Bintang          int    `json:"bintang" binding:"required"`
		BintangLayanan   int    `json:"bintang_layanan"`
		BintangKurir     int    `json:"bintang_kurir"`
		BintangKecepatan int    `json:"bintang_kecepatan"`
		Ulasan           string `json:"ulasan"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default fallback to bintang if not specified
	bintangLayanan := input.BintangLayanan
	if bintangLayanan == 0 {
		bintangLayanan = input.Bintang
	}
	bintangKurir := input.BintangKurir
	if bintangKurir == 0 {
		bintangKurir = input.Bintang
	}
	bintangKecepatan := input.BintangKecepatan
	if bintangKecepatan == 0 {
		bintangKecepatan = input.Bintang
	}

	penilaian := model.Penilaian{
		OrderID:          uint(orderID),
		Bintang:          input.Bintang,
		BintangLayanan:   bintangLayanan,
		BintangKurir:     bintangKurir,
		BintangKecepatan: bintangKecepatan,
		Ulasan:           input.Ulasan,
		TglPenilaian:     time.Now(),
	}

	if err := c.db.Create(&penilaian).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ulasan"})
		return
	}

	// Trigger notification to staff
	go func() {
		var staff []model.User
		if err := c.db.Where("id_role IN (1, 2)").Find(&staff).Error; err == nil {
			namaPelanggan := "Pelanggan"
			if order.Pelanggan.NamaLengkap != "" {
				namaPelanggan = order.Pelanggan.NamaLengkap
			}
			title := "Ulasan Baru ⭐️"
			message := fmt.Sprintf("%s memberikan rating %d bintang untuk pesanan %s.", namaPelanggan, penilaian.Bintang, order.KodeOrder)
			if penilaian.Ulasan != "" {
				message = fmt.Sprintf("%s memberikan rating %d bintang: \"%s\" untuk pesanan %s.", namaPelanggan, penilaian.Bintang, penilaian.Ulasan, order.KodeOrder)
			}
			
			for _, s := range staff {
				notif := model.Notifikasi{
					UserID: s.IDUser,
					Judul:  title,
					Pesan:  message,
					IsRead: false,
				}
				if err := c.db.Create(&notif).Error; err == nil {
					GlobalNotifHub.BroadcastNotification(s.IDUser, notif)
				}
			}
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ulasan berhasil dikirim",
		"data":    penilaian,
	})
}
