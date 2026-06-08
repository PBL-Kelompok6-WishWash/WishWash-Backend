package controller

import (
	"net/http"
	"strconv"

	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
)

type NotifikasiController struct {
	notifikasiRepo repository.NotifikasiRepository
}

func NewNotifikasiController(notifikasiRepo repository.NotifikasiRepository) *NotifikasiController {
	return &NotifikasiController{notifikasiRepo}
}

func (ctrl *NotifikasiController) GetNotifications(c *gin.Context) {
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak terautentikasi"})
		return
	}
	userID := uint(userIDFloat.(float64))

	notifications, err := ctrl.notifikasiRepo.FindAllByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data notifikasi berhasil diambil",
		"data":    notifications,
	})
}

func (ctrl *NotifikasiController) MarkAsRead(c *gin.Context) {
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak terautentikasi"})
		return
	}
	userID := uint(userIDFloat.(float64))

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID notifikasi tidak valid"})
		return
	}

	err = ctrl.notifikasiRepo.MarkAsRead(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifikasi berhasil ditandai telah dibaca",
	})
}

func (ctrl *NotifikasiController) MarkAllAsRead(c *gin.Context) {
	userIDFloat, exists := c.Get("id_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak terautentikasi"})
		return
	}
	userID := uint(userIDFloat.(float64))

	err := ctrl.notifikasiRepo.MarkAllAsRead(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui semua status notifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Semua notifikasi berhasil ditandai telah dibaca",
	})
}
