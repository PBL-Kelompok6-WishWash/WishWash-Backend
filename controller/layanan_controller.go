package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PaketLayananInput struct {
	NamaPaket     string  `json:"nama_paket" binding:"required"`
	DurasiJam     int     `json:"durasi_jam"`
	BiayaTambahan float64 `json:"biaya_tambahan"`
}

type LayananInput struct {
	NamaLayanan     string              `json:"nama_layanan" binding:"required"`
	GambarLayanan   string              `json:"gambar_layanan"`
	JenisSatuan     string              `json:"jenis_satuan" binding:"required"`
	HargaPerSatuan  float64             `json:"harga_per_satuan" binding:"required"`
	ReferensiStatus []string            `json:"referensi_status" binding:"required,min=1"` // Array of status names in sequence
	PaketLayanan    []PaketLayananInput `json:"paket_layanan"`                           // Daftar paket (opsional)
	StatusLayanan    string              `json:"status_layanan"`
	WarnaLayanan     string              `json:"warna_layanan"`
	DeskripsiLayanan string              `json:"deskripsi_layanan"`
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
	
	// Check if each service is used and map to custom map with is_used
	var responseData []gin.H
	for _, l := range layanans {
		isUsed, _ := ctrl.layananRepo.CheckIsUsed(l.IDLayanan)
		responseData = append(responseData, gin.H{
			"id_layanan":        l.IDLayanan,
			"nama_layanan":      l.NamaLayanan,
			"gambar_layanan":    l.GambarLayanan,
			"jenis_satuan":      l.JenisSatuan,
			"harga_per_satuan":  l.HargaPerSatuan,
			"status_layanan":    l.StatusLayanan,
			"referensi_status":  l.ReferensiStatus,
			"paket_layanan":     l.PaketLayanan,
			"warna_layanan":     l.WarnaLayanan,
			"deskripsi_layanan": l.DeskripsiLayanan,
			"is_used":            isUsed,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responseData})
}

func (ctrl *layananController) GetByID(c *gin.Context) {
	id := parseLayananID(c.Param("id"))

	layanan, err := ctrl.layananRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}

	isUsed, _ := ctrl.layananRepo.CheckIsUsed(layanan.IDLayanan)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id_layanan":        layanan.IDLayanan,
			"nama_layanan":      layanan.NamaLayanan,
			"gambar_layanan":    layanan.GambarLayanan,
			"jenis_satuan":      layanan.JenisSatuan,
			"harga_per_satuan":  layanan.HargaPerSatuan,
			"status_layanan":    layanan.StatusLayanan,
			"referensi_status":  layanan.ReferensiStatus,
			"paket_layanan":     layanan.PaketLayanan,
			"warna_layanan":     layanan.WarnaLayanan,
			"deskripsi_layanan": layanan.DeskripsiLayanan,
			"is_used":            isUsed,
		},
	})
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

	// Buat objek Layanan (tanpa gambar dulu agar dapet ID)
	layanan := model.Layanan{
		NamaLayanan:    input.NamaLayanan,
		GambarLayanan:  "",
		JenisSatuan:    input.JenisSatuan,
		HargaPerSatuan: input.HargaPerSatuan,
		StatusLayanan:    input.StatusLayanan,
		WarnaLayanan:     input.WarnaLayanan,
		DeskripsiLayanan: input.DeskripsiLayanan,
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

	var pakets []model.PaketLayanan
	for _, pkt := range input.PaketLayanan {
		pakets = append(pakets, model.PaketLayanan{
			NamaPaket:     strings.TrimSpace(pkt.NamaPaket),
			DurasiJam:     pkt.DurasiJam,
			BiayaTambahan: pkt.BiayaTambahan,
		})
	}
	layanan.PaketLayanan = pakets

	// Simpan
	if err := ctrl.layananRepo.Create(&layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan layanan"})
		return
	}

	// Simpan foto ke subfolder per-entity (sekarang ID sudah ada)
	if input.GambarLayanan != "" {
		entityFolder := utils.BuildEntityFolder(layanan.IDLayanan, input.NamaLayanan)
		gambarPath, err := utils.SaveBase64Image(input.GambarLayanan, "layanan", entityFolder, "layanan_"+input.NamaLayanan)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar layanan tidak valid atau gagal disimpan"})
			return
		}
		layanan.GambarLayanan = gambarPath
		_ = ctrl.layananRepo.Update(&layanan)
	}

	// Ambil data lengkap untuk kembalian
	fullData, _ := ctrl.layananRepo.FindByID(layanan.IDLayanan)
	c.JSON(http.StatusCreated, gin.H{"message": "Layanan berhasil ditambahkan!", "data": fullData})
}

type UpdateLayananInput struct {
	NamaLayanan     *string              `json:"nama_layanan"`
	GambarLayanan   *string              `json:"gambar_layanan"`
	JenisSatuan     *string              `json:"jenis_satuan"`
	HargaPerSatuan  *float64             `json:"harga_per_satuan"`
	ReferensiStatus *[]string            `json:"referensi_status"`
	PaketLayanan    *[]PaketLayananInput `json:"paket_layanan"`
	StatusLayanan    *string              `json:"status_layanan"`
	WarnaLayanan     *string              `json:"warna_layanan"`
	DeskripsiLayanan *string              `json:"deskripsi_layanan"`
}

func (ctrl *layananController) Update(c *gin.Context) {
	id := parseLayananID(c.Param("id"))
	var input UpdateLayananInput

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

	isUsed, _ := ctrl.layananRepo.CheckIsUsed(id)
	if isUsed {
		// Jika digunakan, kita tidak boleh mengubah: nama_layanan, jenis_satuan, referensi_status
		if input.NamaLayanan != nil && *input.NamaLayanan != layanan.NamaLayanan {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama layanan tidak dapat diubah karena telah digunakan dalam transaksi"})
			return
		}
		if input.JenisSatuan != nil && *input.JenisSatuan != layanan.JenisSatuan {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Jenis satuan tidak dapat diubah karena telah digunakan dalam transaksi"})
			return
		}
		if input.ReferensiStatus != nil {
			// Periksa apakah alur status berubah secara jumlah atau nama
			if len(*input.ReferensiStatus) != len(layanan.ReferensiStatus) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Alur status layanan tidak dapat diubah karena telah digunakan dalam transaksi"})
				return
			}
			for idx, sName := range *input.ReferensiStatus {
				if strings.TrimSpace(sName) != layanan.ReferensiStatus[idx].NamaStatus {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Nama alur status layanan tidak dapat diubah karena telah digunakan dalam transaksi"})
					return
				}
			}
		}
	}

	// Update data utama
	namaLayanan := layanan.NamaLayanan
	if input.NamaLayanan != nil {
		namaLayanan = *input.NamaLayanan
		layanan.NamaLayanan = *input.NamaLayanan
	}
	if input.GambarLayanan != nil {
		oldGambarPath := layanan.GambarLayanan
		entityFolder := utils.BuildEntityFolder(layanan.IDLayanan, namaLayanan)

		if *input.GambarLayanan == "" {
			// User sengaja menghapus gambar
			layanan.GambarLayanan = ""
			if oldGambarPath != "" {
				utils.DeleteImageFile(oldGambarPath)
			}
		} else {
			// User mengupload gambar baru (base64)
			gambarPath, err := utils.SaveBase64Image(*input.GambarLayanan, "layanan", entityFolder, "layanan_"+namaLayanan)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan gambar baru"})
				return
			}
			layanan.GambarLayanan = gambarPath

			// Hapus file lama jika path-nya berubah
			if oldGambarPath != "" && oldGambarPath != gambarPath {
				utils.DeleteImageFile(oldGambarPath)
			}
		}
	}
	if input.JenisSatuan != nil {
		layanan.JenisSatuan = *input.JenisSatuan
	}
	if input.HargaPerSatuan != nil {
		layanan.HargaPerSatuan = *input.HargaPerSatuan
	}
	if input.StatusLayanan != nil {
		layanan.StatusLayanan = *input.StatusLayanan
	}
	if input.WarnaLayanan != nil {
		layanan.WarnaLayanan = *input.WarnaLayanan
	}
	if input.DeskripsiLayanan != nil {
		layanan.DeskripsiLayanan = *input.DeskripsiLayanan
	}

	if err := ctrl.layananRepo.Update(layanan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate layanan"})
		return
	}

	// Update data status (hanya jika dikirim)
	if input.ReferensiStatus != nil {
		var statuses []model.ReferensiStatusLayanan
		for i, statusName := range *input.ReferensiStatus {
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
	}

	// Update data paket (hanya jika dikirim)
	if input.PaketLayanan != nil {
		var pakets []model.PaketLayanan
		for _, pkt := range *input.PaketLayanan {
			pakets = append(pakets, model.PaketLayanan{
				LayananID:     id,
				NamaPaket:     strings.TrimSpace(pkt.NamaPaket),
				DurasiJam:     pkt.DurasiJam,
				BiayaTambahan: pkt.BiayaTambahan,
			})
		}
		if err := ctrl.layananRepo.UpdatePaketLayanan(id, pakets); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate paket layanan"})
			return
		}
	}

	fullData, _ := ctrl.layananRepo.FindByID(id)
	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil diperbarui!", "data": fullData})
}

func (ctrl *layananController) Delete(c *gin.Context) {
	id := parseLayananID(c.Param("id"))

	layanan, err := ctrl.layananRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan tidak ditemukan"})
		return
	}

	isUsed, _ := ctrl.layananRepo.CheckIsUsed(id)
	if isUsed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Layanan tidak dapat dihapus karena telah digunakan dalam transaksi"})
		return
	}

	// Hapus seluruh folder entity dari disk (bersih sekaligus)
	entityFolder := utils.BuildEntityFolder(layanan.IDLayanan, layanan.NamaLayanan)
	utils.DeleteImageFolder("layanan", entityFolder)

	if err := ctrl.layananRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Layanan berhasil dihapus"})
}
