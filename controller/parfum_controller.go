package controller

import (
	"net/http"
	"strconv"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
)

type ParfumController struct {
	parfumRepo repository.ParfumRepository
}

func NewParfumController(parfumRepo repository.ParfumRepository) *ParfumController {
	return &ParfumController{parfumRepo}
}

func (c *ParfumController) GetAll(ctx *gin.Context) {
	parfums, err := c.parfumRepo.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data parfum"})
		return
	}

	var responseData []gin.H
	for _, p := range parfums {
		isUsed, _ := c.parfumRepo.CheckIsUsed(int(p.IDParfum))
		responseData = append(responseData, gin.H{
			"id_parfum":     p.IDParfum,
			"nama_parfum":   p.NamaParfum,
			"keterangan":    p.Keterangan,
			"status_parfum": p.StatusParfum,
			"is_used":       isUsed,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data Parfum Berhasil Diambil",
		"data":    responseData,
	})
}

func (c *ParfumController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parfum tidak valid"})
		return
	}

	parfum, err := c.parfumRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data parfum tidak ditemukan"})
		return
	}

	isUsed, _ := c.parfumRepo.CheckIsUsed(id)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data Parfum Berhasil Diambil",
		"data": gin.H{
			"id_parfum":     parfum.IDParfum,
			"nama_parfum":   parfum.NamaParfum,
			"keterangan":    parfum.Keterangan,
			"status_parfum": parfum.StatusParfum,
			"is_used":       isUsed,
		},
	})
}

func (c *ParfumController) Create(ctx *gin.Context) {
	var input struct {
		NamaParfum string `json:"nama_parfum" binding:"required"`
		Keterangan string `json:"keterangan"`
		Status     string `json:"status_parfum"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parfum := model.Parfum{
		NamaParfum:   input.NamaParfum,
		Keterangan:   input.Keterangan,
		StatusParfum: input.Status,
	}

	newParfum, err := c.parfumRepo.Create(parfum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data parfum"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data Parfum Berhasil Disimpan",
		"data":    newParfum,
	})
}

func (c *ParfumController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parfum tidak valid"})
		return
	}

	var input struct {
		NamaParfum *string `json:"nama_parfum"`
		Keterangan *string `json:"keterangan"`
		Status     *string `json:"status_parfum"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parfum, err := c.parfumRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data parfum tidak ditemukan"})
		return
	}

	isUsed, _ := c.parfumRepo.CheckIsUsed(id)
	if isUsed {
		if input.NamaParfum != nil && *input.NamaParfum != parfum.NamaParfum {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nama parfum tidak dapat diubah karena telah digunakan dalam transaksi"})
			return
		}
	}

	if input.NamaParfum != nil {
		parfum.NamaParfum = *input.NamaParfum
	}
	if input.Keterangan != nil {
		parfum.Keterangan = *input.Keterangan
	}
	if input.Status != nil {
		parfum.StatusParfum = *input.Status
	}

	updatedParfum, err := c.parfumRepo.Update(parfum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate data parfum"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data Parfum Berhasil Diupdate",
		"data":    updatedParfum,
	})
}

func (c *ParfumController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parfum tidak valid"})
		return
	}

	parfum, err := c.parfumRepo.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data parfum tidak ditemukan"})
		return
	}

	isUsed, _ := c.parfumRepo.CheckIsUsed(id)
	if isUsed {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parfum tidak dapat dihapus karena telah digunakan dalam transaksi"})
		return
	}

	err = c.parfumRepo.Delete(parfum)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data parfum"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data Parfum Berhasil Dihapus",
	})
}
