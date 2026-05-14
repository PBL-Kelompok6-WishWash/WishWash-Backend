package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
)

type PromoInput struct {
	KodePromo        string  `json:"kode_promo" binding:"required"`
	NamaPromo        string  `json:"nama_promo" binding:"required"`
	Deskripsi        string  `json:"deskripsi"`
	TipePromo        string  `json:"tipe_promo" binding:"required"` // "Persentase" atau "Nominal"
	NominalPotongan  float64 `json:"nominal_potongan" binding:"required"`
	MinimalOrder     float64 `json:"minimal_order"`
	MaksimalPotongan float64 `json:"maksimal_potongan"`
	TglMulai         string  `json:"tgl_mulai" binding:"required"`   // format: "2006-01-02"
	TglBerakhir      string  `json:"tgl_berakhir" binding:"required"` // format: "2006-01-02"
	StatusPromo      string  `json:"status_promo" binding:"required"` // "Aktif" atau "Tidak Aktif"
	GambarPromo      string  `json:"gambar_promo"`
}

type PromoController struct {
	promoRepo repository.PromoRepository
}

func NewPromoController(promoRepo repository.PromoRepository) *PromoController {
	return &PromoController{promoRepo}
}

func (c *PromoController) GetAll(ctx *gin.Context) {
	promos, err := c.promoRepo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data promo"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data Promo Berhasil Diambil", "data": promos})
}

func (c *PromoController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID promo tidak valid"})
		return
	}

	promo, err := c.promoRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data promo tidak ditemukan"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data Promo Berhasil Diambil", "data": promo})
}

func (c *PromoController) Create(ctx *gin.Context) {
	var input PromoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tglMulai, err := time.Parse("2006-01-02", input.TglMulai)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal mulai tidak valid, gunakan YYYY-MM-DD"})
		return
	}
	tglBerakhir, err := time.Parse("2006-01-02", input.TglBerakhir)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal berakhir tidak valid, gunakan YYYY-MM-DD"})
		return
	}

	promo := model.Promo{
		KodePromo:        input.KodePromo,
		NamaPromo:        input.NamaPromo,
		Deskripsi:        input.Deskripsi,
		TipePromo:        input.TipePromo,
		NominalPotongan:  input.NominalPotongan,
		MinimalOrder:     input.MinimalOrder,
		MaksimalPotongan: input.MaksimalPotongan,
		TglMulai:         tglMulai,
		TglBerakhir:      tglBerakhir,
		StatusPromo:      input.StatusPromo,
		GambarPromo:      input.GambarPromo,
	}

	newPromo, err := c.promoRepo.Create(promo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data promo"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Data Promo Berhasil Disimpan", "data": newPromo})
}

func (c *PromoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID promo tidak valid"})
		return
	}

	var input PromoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promo, err := c.promoRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data promo tidak ditemukan"})
		return
	}

	tglMulai, err := time.Parse("2006-01-02", input.TglMulai)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal mulai tidak valid, gunakan YYYY-MM-DD"})
		return
	}
	tglBerakhir, err := time.Parse("2006-01-02", input.TglBerakhir)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal berakhir tidak valid, gunakan YYYY-MM-DD"})
		return
	}

	promo.KodePromo        = input.KodePromo
	promo.NamaPromo        = input.NamaPromo
	promo.Deskripsi        = input.Deskripsi
	promo.TipePromo        = input.TipePromo
	promo.NominalPotongan  = input.NominalPotongan
	promo.MinimalOrder     = input.MinimalOrder
	promo.MaksimalPotongan = input.MaksimalPotongan
	promo.TglMulai         = tglMulai
	promo.TglBerakhir      = tglBerakhir
	promo.StatusPromo      = input.StatusPromo
	promo.GambarPromo      = input.GambarPromo

	updated, err := c.promoRepo.Update(promo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data promo"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data Promo Berhasil Diupdate", "data": updated})
}

func (c *PromoController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID promo tidak valid"})
		return
	}

	promo, err := c.promoRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data promo tidak ditemukan"})
		return
	}

	if err := c.promoRepo.Delete(promo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data promo"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data Promo Berhasil Dihapus"})
}
