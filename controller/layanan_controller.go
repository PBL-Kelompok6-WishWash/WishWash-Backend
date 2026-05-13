package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LayananInput struct {
	NamaLayanan     string   `json:"nama_layanan" binding:"required"`
	GambarLayanan   string   `json:"gambar_layanan"`
	JenisSatuan     string   `json:"jenis_satuan" binding:"required"`
	HargaPerSatuan  float64  `json:"harga_per_satuan" binding:"required"`
	ReferensiStatus []string `json:"referensi_status" binding:"required,min=1"` // Array of status names in sequence
}

type LayananController interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type layananController struct {
	layananRepo repository.LayananRepository
}

func NewLayananController(lRepo repository.LayananRepository) LayananController {
	return &layananController{lRepo}
}

func parseLayananID(idStr string) uint {
	id, _ := strconv.ParseUint(idStr, 10, 32)
	return uint(id)
}

func (ctrl *layananController) GetAll(c *gin.Context) {
	layanans, err := ctrl.layananRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": layanans})
}

func (ctrl *layananController) GetByID(c *gin.Context) {
	id := parseLayananID(c.Param("id"))

	layanan, err := ctrl.layananRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": layanan})
}

func (ctrl *layananController) Create(c *gin.Context) {
	var input LayananInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errorMsgs []string
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				switch fieldErr.Field() {
				case "NamaLayanan":
					errorMsgs = append(errorMsgs, "Nama Layanan tidak boleh kosong")
				case "JenisSatuan":
					errorMsgs = append(errorMsgs, "Jenis Satuan tidak boleh kosong")
				case "HargaPerSatuan":
					errorMsgs = append(errorMsgs, "Harga tidak valid")
				case "ReferensiStatus":
					errorMsgs = append(errorMsgs, "Minimal harus ada 1 status layanan")
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMsgs, ", ")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Buat objek Layanan
	layanan := model.Layanan{
		NamaLayanan:    input.NamaLayanan,
		GambarLayanan:  input.GambarLayanan,
		JenisSatuan:    input.JenisSatuan,
		HargaPerSatuan: input.HargaPerSatuan,
	}

	// Konversi array string menjadi array struct ReferensiStatusLayanan
	var statuses []model.ReferensiStatusLayanan
	for i, statusName := range input.ReferensiStatus {
		statuses = append(statuses, model.ReferensiStatusLayanan{
			NamaStatus:  strings.TrimSpace(statusName),
			UrutanTahap: i + 1, // Urutan mulai dari 1
		})
	}
	layanan.ReferensiStatus = statuses

	// Simpan
	if err := ctrl.layananRepo.Create(&layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan layanan"})
		return
	}

	// Ambil data lengkap untuk kembalian
	fullData, _ := ctrl.layananRepo.FindByID(layanan.IDLayanan)
	c.JSON(http.StatusCreated, gin.H{"message": "Layanan berhasil ditambahkan!", "data": fullData})
}

func (ctrl *layananController) Update(c *gin.Context) {
	id := parseLayananID(c.Param("id"))
	var input LayananInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Cek apakah layanan ada
	layanan, err := ctrl.layananRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}

	// Update data utama
	layanan.NamaLayanan = input.NamaLayanan
	layanan.GambarLayanan = input.GambarLayanan
	layanan.JenisSatuan = input.JenisSatuan
	layanan.HargaPerSatuan = input.HargaPerSatuan

	if err := ctrl.layananRepo.Update(layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate layanan"})
		return
	}

	// Update data status
	var statuses []model.ReferensiStatusLayanan
	for i, statusName := range input.ReferensiStatus {
		statuses = append(statuses, model.ReferensiStatusLayanan{
			LayananID:   id,
			NamaStatus:  strings.TrimSpace(statusName),
			UrutanTahap: i + 1,
		})
	}

	if err := ctrl.layananRepo.UpdateStatusLayanan(id, statuses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate referensi status layanan"})
		return
	}

	fullData, _ := ctrl.layananRepo.FindByID(id)
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil diperbarui!", "data": fullData})
}

func (ctrl *layananController) Delete(c *gin.Context) {
	id := parseLayananID(c.Param("id"))

	_, err := ctrl.layananRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}

	if err := ctrl.layananRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil dihapus"})
}
