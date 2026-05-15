package controller

import (
	"net/http"
	// "strconv"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MetodePembayaranController struct {
	DB *gorm.DB
}

func NewMetodePembayaranController(db *gorm.DB) *MetodePembayaranController {
	return &MetodePembayaranController{DB: db}
}

// GetAll godoc
func (ctrl *MetodePembayaranController) GetAll(c *gin.Context) {
	var mps []model.MetodePembayaran
	if err := ctrl.DB.Order("id_metode_pembayaran asc").Find(&mps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data metode pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mps})
}

// GetByID godoc
func (ctrl *MetodePembayaranController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var mp model.MetodePembayaran
	if err := ctrl.DB.First(&mp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metode pembayaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mp})
}

// Create godoc
func (ctrl *MetodePembayaranController) Create(c *gin.Context) {
	var mp model.MetodePembayaran
	if err := c.ShouldBindJSON(&mp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&mp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat metode pembayaran"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": mp})
}

// Update godoc
func (ctrl *MetodePembayaranController) Update(c *gin.Context) {
	id := c.Param("id")
	var mp model.MetodePembayaran
	if err := ctrl.DB.First(&mp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metode pembayaran tidak ditemukan"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Model(&mp).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui metode pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": mp})
}

// Delete godoc
func (ctrl *MetodePembayaranController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&model.MetodePembayaran{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus metode pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Metode pembayaran berhasil dihapus"})
}
