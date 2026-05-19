package controller

import (
	"net/http"
	"strconv"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
)

type AlamatController interface {
	GetAlamatPelanggan(c *gin.Context)
	CreateAlamat(c *gin.Context)
	UpdateAlamat(c *gin.Context)
	SetPrimaryAlamat(c *gin.Context)
	DeleteAlamat(c *gin.Context)
}

type alamatController struct {
	alamatRepo    repository.AlamatRepository
	pelangganRepo repository.PelangganRepository
}

func NewAlamatController(aRepo repository.AlamatRepository, pRepo repository.PelangganRepository) AlamatController {
	return &alamatController{aRepo, pRepo}
}

func (ctrl *alamatController) getPelangganIDFromContext(c *gin.Context) (uint, error) {
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

func (ctrl *alamatController) GetAlamatPelanggan(c *gin.Context) {
	pelangganID, err := ctrl.getPelangganIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	alamats, err := ctrl.alamatRepo.FindAllByPelangganID(pelangganID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data alamat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data alamat berhasil diambil",
		"data":    alamats,
	})
}

func (ctrl *alamatController) CreateAlamat(c *gin.Context) {
	pelangganID, err := ctrl.getPelangganIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	var input struct {
		AlamatLengkap string `json:"alamat_lengkap" binding:"required"`
		TipeAlamat    string `json:"tipe_alamat" binding:"required"`
		NamaPenerima  string `json:"nama_penerima" binding:"required"`
		NoHpPenerima  string `json:"nohp_penerima" binding:"required"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	alamat := model.Alamat{
		PelangganID:   pelangganID,
		AlamatLengkap: input.AlamatLengkap,
		TipeAlamat:    input.TipeAlamat,
		NamaPenerima:  input.NamaPenerima,
		NoHpPenerima:  input.NoHpPenerima,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
	}

	if err := ctrl.alamatRepo.Create(&alamat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan alamat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alamat berhasil ditambahkan",
		"data":    alamat,
	})
}

func (ctrl *alamatController) UpdateAlamat(c *gin.Context) {
	pelangganID, err := ctrl.getPelangganIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	idAlamatParam := c.Param("id")
	idAlamat, err := strconv.ParseUint(idAlamatParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Alamat tidak valid"})
		return
	}

	var input struct {
		AlamatLengkap string `json:"alamat_lengkap" binding:"required"`
		TipeAlamat    string `json:"tipe_alamat" binding:"required"`
		NamaPenerima  string `json:"nama_penerima" binding:"required"`
		NoHpPenerima  string `json:"nohp_penerima" binding:"required"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	alamat := model.Alamat{
		IDAlamat:      uint(idAlamat),
		PelangganID:   pelangganID,
		AlamatLengkap: input.AlamatLengkap,
		TipeAlamat:    input.TipeAlamat,
		NamaPenerima:  input.NamaPenerima,
		NoHpPenerima:  input.NoHpPenerima,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
	}

	if err := ctrl.alamatRepo.Update(&alamat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah alamat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alamat berhasil diubah",
		"data":    alamat,
	})
}

func (ctrl *alamatController) SetPrimaryAlamat(c *gin.Context) {
	pelangganID, err := ctrl.getPelangganIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	idAlamatParam := c.Param("id")
	idAlamat, err := strconv.ParseUint(idAlamatParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Alamat tidak valid"})
		return
	}

	if err := ctrl.alamatRepo.SetPrimary(uint(idAlamat), pelangganID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah alamat utama"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alamat utama berhasil diubah",
	})
}

func (ctrl *alamatController) DeleteAlamat(c *gin.Context) {
	pelangganID, err := ctrl.getPelangganIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}

	idAlamatParam := c.Param("id")
	idAlamat, err := strconv.ParseUint(idAlamatParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Alamat tidak valid"})
		return
	}

	if err := ctrl.alamatRepo.Delete(uint(idAlamat), pelangganID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus alamat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alamat berhasil dihapus",
	})
}
