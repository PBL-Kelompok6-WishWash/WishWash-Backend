package controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type NotifikasiHub struct {
	Clients map[uint][]*websocket.Conn
	Mu      sync.Mutex
}

var GlobalNotifHub = &NotifikasiHub{
	Clients: make(map[uint][]*websocket.Conn),
}

func (hub *NotifikasiHub) BroadcastNotification(userID uint, notif model.Notifikasi) {
	hub.Mu.Lock()
	defer hub.Mu.Unlock()
	conns, exists := hub.Clients[userID]
	if !exists {
		log.Printf("🔔 [NotifHub] No active WS connections found for UserID %d. Notification saved but not broadcasted.", userID)
		return
	}
	log.Printf("🔔 [NotifHub] Broadcasting notification to UserID %d: %s. Connections count: %d", userID, notif.Judul, len(conns))
	var activeConns []*websocket.Conn
	for _, conn := range conns {
		err := conn.WriteJSON(gin.H{
			"id_notifikasi": notif.IDNotifikasi,
			"id_user":       notif.UserID,
			"judul":         notif.Judul,
			"pesan":         notif.Pesan,
			"is_read":       notif.IsRead,
			"created_at":    notif.CreatedAt.Format(time.RFC3339),
		})
		if err == nil {
			activeConns = append(activeConns, conn)
		} else {
			log.Printf("🔴 [NotifHub] Failed to write JSON to WS for UserID %d: %v", userID, err)
			conn.Close()
		}
	}
	if len(activeConns) > 0 {
		hub.Clients[userID] = activeConns
	} else {
		delete(hub.Clients, userID)
	}
}

var notifUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotifikasiController struct {
	notifikasiRepo repository.NotifikasiRepository
}

func NewNotifikasiController(notifikasiRepo repository.NotifikasiRepository) *NotifikasiController {
	return &NotifikasiController{notifikasiRepo}
}

func (ctrl *NotifikasiController) HandleNotifWS(c *gin.Context) {
	userIDStr := c.Query("id_user")
	userIDVal, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID User tidak valid"})
		return
	}
	userID := uint(userIDVal)

	conn, err := notifUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("🔴 [NotifWS] WS Upgrade error for UserID %d: %v", userID, err)
		return
	}
	defer conn.Close()

	log.Printf("🔔 [NotifWS] Client connected to Notification WS: UserID %d", userID)

	GlobalNotifHub.Mu.Lock()
	GlobalNotifHub.Clients[userID] = append(GlobalNotifHub.Clients[userID], conn)
	GlobalNotifHub.Mu.Unlock()

	defer func() {
		log.Printf("🔔 [NotifWS] Client disconnected from Notification WS: UserID %d", userID)
		GlobalNotifHub.Mu.Lock()
		conns := GlobalNotifHub.Clients[userID]
		for i, v := range conns {
			if v == conn {
				GlobalNotifHub.Clients[userID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(GlobalNotifHub.Clients[userID]) == 0 {
			delete(GlobalNotifHub.Clients, userID)
		}
		GlobalNotifHub.Mu.Unlock()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
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
