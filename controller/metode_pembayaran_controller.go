package controller

import (
	"net/http"
	// "strconv"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/utils"
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

	// 1. Simpan detail tanpa gambarMetode dulu untuk dapat ID
	gambarBase64 := mp.GambarMetode
	mp.GambarMetode = ""

	if err := ctrl.DB.Create(&mp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat metode pembayaran"})
		return
	}

	// 2. Simpan gambar ke subfolder per-entity menggunakan ID
	if gambarBase64 != "" {
		entityFolder := utils.BuildEntityFolder(mp.IDMetodePembayaran, mp.NamaMetode)
		gambarPath, err := utils.SaveBase64Image(gambarBase64, "metode_bayar", entityFolder, "metode_"+mp.NamaMetode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar metode pembayaran tidak valid atau gagal disimpan"})
			return
		}
		mp.GambarMetode = gambarPath
		_ = ctrl.DB.Model(&mp).Update("gambar_metode", gambarPath)
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

	if gambarVal, exists := input["gambar_metode"]; exists {
		gambarStr, ok := gambarVal.(string)
		if ok {
			oldGambarPath := mp.GambarMetode
			namaMetode := mp.NamaMetode
			if namaVal, nameExists := input["nama_metode"]; nameExists {
				if nameStr, nameOk := namaVal.(string); nameOk && nameStr != "" {
					namaMetode = nameStr
				}
			}

			if gambarStr == "" {
				// User klik X lalu simpan (hapus gambar)
				input["gambar_metode"] = ""
				if oldGambarPath != "" {
					utils.DeleteImageFile(oldGambarPath)
				}
			} else {
				// User upload gambar baru
				entityFolder := utils.BuildEntityFolder(mp.IDMetodePembayaran, namaMetode)
				gambarPath, err := utils.SaveBase64Image(gambarStr, "metode_bayar", entityFolder, "metode_"+namaMetode)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menyimpan gambar baru"})
					return
				}
				input["gambar_metode"] = gambarPath

				// Hapus file lama jika path-nya berubah
				if oldGambarPath != "" && oldGambarPath != gambarPath {
					utils.DeleteImageFile(oldGambarPath)
				}
			}
		}
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
	var mp model.MetodePembayaran
	if err := ctrl.DB.First(&mp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Metode pembayaran tidak ditemukan"})
		return
	}

	// Hapus seluruh folder entity dari disk (bersih sekaligus)
	entityFolder := utils.BuildEntityFolder(mp.IDMetodePembayaran, mp.NamaMetode)
	utils.DeleteImageFolder("metode_bayar", entityFolder)

	if err := ctrl.DB.Delete(&mp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus metode pembayaran"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Metode pembayaran berhasil dihapus"})
}
