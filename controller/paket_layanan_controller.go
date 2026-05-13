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

type PaketLayananInput struct {
	NamaPaket     string  `json:"nama_paket" binding:"required"`
	DurasiJam     int     `json:"durasi_jam" binding:"required"`
	BiayaTambahan float64 `json:"biaya_tambahan"`
}

type PaketLayananController interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type paketLayananController struct {
	paketRepo repository.PaketLayananRepository
}

func NewPaketLayananController(pRepo repository.PaketLayananRepository) PaketLayananController {
	return &paketLayananController{pRepo}
}

func parsePaketID(idStr string) uint {
	id, _ := strconv.ParseUint(idStr, 10, 32)
	return uint(id)
}

func (ctrl *paketLayananController) GetAll(c *gin.Context) {
	pakets, err := ctrl.paketRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data paket layanan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pakets})
}

func (ctrl *paketLayananController) GetByID(c *gin.Context) {
	id := parsePaketID(c.Param("id"))

	paket, err := ctrl.paketRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paket layanan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paket})
}

func (ctrl *paketLayananController) Create(c *gin.Context) {
	var input PaketLayananInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errorMsgs []string
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrs {
				switch fieldErr.Field() {
				case "NamaPaket":
					errorMsgs = append(errorMsgs, "Nama Paket tidak boleh kosong")
				case "DurasiJam":
					errorMsgs = append(errorMsgs, "Durasi Jam tidak boleh kosong")
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": strings.Join(errorMsgs, ", ")})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	paket := model.PaketLayanan{
		NamaPaket:     input.NamaPaket,
		DurasiJam:     input.DurasiJam,
		BiayaTambahan: input.BiayaTambahan,
	}

	if err := ctrl.paketRepo.Create(&paket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan paket layanan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Paket layanan berhasil ditambahkan!", "data": paket})
}

func (ctrl *paketLayananController) Update(c *gin.Context) {
	id := parsePaketID(c.Param("id"))
	var input PaketLayananInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	paket, err := ctrl.paketRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paket layanan tidak ditemukan"})
		return
	}

	paket.NamaPaket = input.NamaPaket
	paket.DurasiJam = input.DurasiJam
	paket.BiayaTambahan = input.BiayaTambahan

	if err := ctrl.paketRepo.Update(paket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate paket layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paket layanan berhasil diperbarui!", "data": paket})
}

func (ctrl *paketLayananController) Delete(c *gin.Context) {
	id := parsePaketID(c.Param("id"))

	_, err := ctrl.paketRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paket layanan tidak ditemukan"})
		return
	}

	if err := ctrl.paketRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus paket layanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Paket layanan berhasil dihapus"})
}
