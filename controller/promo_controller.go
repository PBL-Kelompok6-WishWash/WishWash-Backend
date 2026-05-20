package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/utils"
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
		GambarPromo:      "", // Simpan kosong dulu
	}

	newPromo, err := c.promoRepo.Create(promo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data promo"})
		return
	}

	// Setelah dapet ID, simpan gambar ke subfolder
	if input.GambarPromo != "" {
		entityFolder := utils.BuildEntityFolder(newPromo.IDPromo, newPromo.NamaPromo)
		gambarPath, err := utils.SaveBase64Image(input.GambarPromo, "promo", entityFolder, "promo_"+input.KodePromo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar promo tidak valid atau gagal disimpan"})
			return
		}
		newPromo.GambarPromo = gambarPath
		newPromo, _ = c.promoRepo.Update(newPromo)
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Data Promo Berhasil Disimpan", "data": newPromo})
}

type UpdatePromoInput struct {
	KodePromo        *string  `json:"kode_promo"`
	NamaPromo        *string  `json:"nama_promo"`
	Deskripsi        *string  `json:"deskripsi"`
	TipePromo        *string  `json:"tipe_promo"`
	NominalPotongan  *float64 `json:"nominal_potongan"`
	MinimalOrder     *float64 `json:"minimal_order"`
	MaksimalPotongan *float64 `json:"maksimal_potongan"`
	TglMulai         *string  `json:"tgl_mulai"`
	TglBerakhir      *string  `json:"tgl_berakhir"`
	StatusPromo      *string  `json:"status_promo"`
	GambarPromo      *string  `json:"gambar_promo"`
}

func (c *PromoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID promo tidak valid"})
		return
	}

	var input UpdatePromoInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	promo, err := c.promoRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data promo tidak ditemukan"})
		return
	}

	kodePromo := promo.KodePromo
	if input.KodePromo != nil {
		kodePromo = *input.KodePromo
		promo.KodePromo = *input.KodePromo
	}
	if input.NamaPromo != nil {
		promo.NamaPromo = *input.NamaPromo
	}
	if input.Deskripsi != nil {
		promo.Deskripsi = *input.Deskripsi
	}
	if input.TipePromo != nil {
		promo.TipePromo = *input.TipePromo
	}
	if input.NominalPotongan != nil {
		promo.NominalPotongan = *input.NominalPotongan
	}
	if input.MinimalOrder != nil {
		promo.MinimalOrder = *input.MinimalOrder
	}
	if input.MaksimalPotongan != nil {
		promo.MaksimalPotongan = *input.MaksimalPotongan
	}
	if input.TglMulai != nil {
		tglMulai, err := time.Parse("2006-01-02", *input.TglMulai)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal mulai tidak valid, gunakan YYYY-MM-DD"})
			return
		}
		promo.TglMulai = tglMulai
	}
	if input.TglBerakhir != nil {
		tglBerakhir, err := time.Parse("2006-01-02", *input.TglBerakhir)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal berakhir tidak valid, gunakan YYYY-MM-DD"})
			return
		}
		promo.TglBerakhir = tglBerakhir
	}
	if input.StatusPromo != nil {
		promo.StatusPromo = *input.StatusPromo
	}
	if input.GambarPromo != nil {
		oldGambarPath := promo.GambarPromo
		entityFolder := utils.BuildEntityFolder(promo.IDPromo, promo.NamaPromo)

		if *input.GambarPromo == "" {
			// User klik X lalu simpan (hapus gambar)
			promo.GambarPromo = ""
			if oldGambarPath != "" {
				utils.DeleteImageFile(oldGambarPath)
			}
		} else {
			// User upload gambar baru
			gambarPath, err := utils.SaveBase64Image(*input.GambarPromo, "promo", entityFolder, "promo_"+kodePromo)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan gambar promo baru"})
				return
			}
			promo.GambarPromo = gambarPath

			// Hapus file lama jika path-nya berubah
			if oldGambarPath != "" && oldGambarPath != gambarPath {
				utils.DeleteImageFile(oldGambarPath)
			}
		}
	}

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

	// Hapus seluruh folder entity dari disk (bersih sekaligus)
	entityFolder := utils.BuildEntityFolder(promo.IDPromo, promo.NamaPromo)
	utils.DeleteImageFolder("promo", entityFolder)

	if err := c.promoRepo.Delete(promo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data promo"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data Promo Berhasil Dihapus"})
}
