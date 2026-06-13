package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/PBL-Kelompok6-WishWash/backend/config"
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/PBL-Kelompok6-WishWash/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader untuk mengubah koneksi HTTP biasa menjadi WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Mengizinkan koneksi dari hp/flutter tanpa kendala CORS
	},
}

// Struktur untuk menyimpan koneksi aktif di dalam setiap Room Chat
type RoomPool struct {
	Clients map[*websocket.Conn]uint // Menyimpan pointer koneksi dan ID User-nya
	Mu      sync.Mutex
}

// Map global untuk menampung semua room chat yang sedang aktif mengobrol
var activeRooms = make(map[string]*RoomPool)
var roomsMu sync.Mutex

type ChatController struct {
	repo      repository.ChatRepository
	notifRepo repository.NotifikasiRepository
}

func NewChatController(repo repository.ChatRepository, notifRepo repository.NotifikasiRepository) *ChatController {
	return &ChatController{repo: repo, notifRepo: notifRepo}
}

// 1. HTTP Endpoint: Mengambil semua riwayat chat lama
func (c *ChatController) GetMessages(ctx *gin.Context) {
	roomIDStr := ctx.Param("id_room_chat")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID Room tidak valid"})
		return
	}

	userIDFloat, exists := ctx.Get("id_user")
	var userID uint
	if exists {
		userID = uint(userIDFloat.(float64))
	}

	if userID > 0 {
		config.DB.Model(&model.PesanChat{}).
			Where("id_room_chat = ? AND id_user != ?", uint(roomID), userID).
			Update("status_baca", true)
	}

	messages, err := c.repo.GetMessagesByRoomID(uint(roomID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat pesan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": messages})
}

// 1.5. HTTP Endpoint: Mengambil daftar Room Chat aktif untuk user yang login
func (c *ChatController) GetRooms(ctx *gin.Context) {
	userIDFloat, exists := ctx.Get("id_user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	rooms, err := c.repo.GetRoomsByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar room chat"})
		return
	}

	type RoomResponse struct {
		IDRoomChat  uint            `json:"id_room_chat"`
		OrderID     uint            `json:"id_order"`
		WaktuDibuat time.Time       `json:"waktu_dibuat"`
		Order       model.Order     `json:"Order"`
		LastMessage *model.PesanChat `json:"LastMessage"`
		UnreadCount int64           `json:"unread_count"`
	}

	var responseData []RoomResponse
	for _, room := range rooms {
		var lastMsg model.PesanChat
		errLast := config.DB.Where("id_room_chat = ?", room.IDRoomChat).
			Order("waktu_kirim desc").
			Preload("ChatGambar").
			First(&lastMsg).Error
		
		var lastMsgPtr *model.PesanChat
		if errLast == nil {
			if len(lastMsg.ChatGambar) > 0 {
				lastMsg.PathGambar = lastMsg.ChatGambar[0].PathGambar
			}
			lastMsgPtr = &lastMsg
		}
		
		var unreadCount int64
		config.DB.Model(&model.PesanChat{}).
			Where("id_room_chat = ? AND status_baca = ? AND id_user != ?", room.IDRoomChat, false, userID).
			Count(&unreadCount)
		
		responseData = append(responseData, RoomResponse{
			IDRoomChat:  room.IDRoomChat,
			OrderID:     room.OrderID,
			WaktuDibuat: room.WaktuDibuat,
			Order:       room.Order,
			LastMessage: lastMsgPtr,
			UnreadCount: unreadCount,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

// 1.6. HTTP Endpoint: Mendapatkan atau membuat Room Chat baru berdasarkan ID Order
func (c *ChatController) GetOrCreateRoom(ctx *gin.Context) {
	orderIDStr := ctx.Param("id_order")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID Order tidak valid"})
		return
	}

	room, err := c.repo.GetOrCreateRoomByOrderID(uint(orderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan atau membuat room chat"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": room})
}

// 2. WebSocket Endpoint: Menangani jabat tangan (Handshake) dan lempar-tangkap pesan instan
func (c *ChatController) HandleWS(ctx *gin.Context) {
	roomIDStr := ctx.Param("id_room_chat")
	
	// Mengambil ID User dari query parameter (misal: ws://.../ws?id_user=5)
	userIDStr := ctx.Query("id_user")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	// Upgrade koneksi HTTP ke WebSocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Kelola ruangan obrolan agar pesan tidak nyasar ke room lain
	roomsMu.Lock()
	if activeRooms[roomIDStr] == nil {
		activeRooms[roomIDStr] = &RoomPool{
			Clients: make(map[*websocket.Conn]uint),
		}
	}
	pool := activeRooms[roomIDStr]
	roomsMu.Unlock()

	// Daftarkan koneksi baru ke dalam pool ruangan
	pool.Mu.Lock()
	pool.Clients[conn] = uint(userID)

	// Cek apakah ada user lain di pool yang online
	otherOnline := false
	for _, uid := range pool.Clients {
		if uid != uint(userID) {
			otherOnline = true
			break
		}
	}

	// Kirim status awal ke user yang baru terhubung
	_ = conn.WriteJSON(gin.H{
		"type":   "initial_status",
		"online": otherOnline,
	})

	// Broadcast status online ke user lain di room yang sama
	for client := range pool.Clients {
		if pool.Clients[client] != uint(userID) {
			_ = client.WriteJSON(gin.H{
				"type":      "status",
				"id_user":   userID,
				"online":    true,
			})
		}
	}
	pool.Mu.Unlock()

	// Bersihkan koneksi jika user menutup aplikasi atau keluar dari halaman chat
	defer func() {
		pool.Mu.Lock()
		delete(pool.Clients, conn)

		// Broadcast status offline ke user lain di room yang sama
		for client := range pool.Clients {
			_ = client.WriteJSON(gin.H{
				"type":      "status",
				"id_user":   userID,
				"online":    false,
			})
		}
		pool.Mu.Unlock()
	}()

	// Loop terus-menerus untuk mendengarkan apakah ada pesan masuk dari Flutter
	for {
		var incoming struct {
			TeksPesan    string `json:"teks_pesan"`
			Base64Gambar string `json:"base64_gambar"`
			IsTyping     *bool  `json:"is_typing"`
		}

		// Membaca pesan JSON dari Flutter
		err := conn.ReadJSON(&incoming)
		if err != nil {
			break // Keluar dari loop jika koneksi terputus
		}

		// Jika pesan berisi status typing, lakukan broadcast tanpa simpan ke DB
		if incoming.IsTyping != nil {
			pool.Mu.Lock()
			for client := range pool.Clients {
				if pool.Clients[client] != uint(userID) {
					_ = client.WriteJSON(gin.H{
						"type":      "typing",
						"id_user":   userID,
						"is_typing": *incoming.IsTyping,
					})
				}
			}
			pool.Mu.Unlock()
			continue
		}

		roomIDUint, _ := strconv.ParseUint(roomIDStr, 10, 32)

		// Bungkus ke dalam model GORM untuk disimpan ke database
		msg := model.PesanChat{
			RoomChatID: uint(roomIDUint),
			UserID:     uint(userID),
			TeksPesan:  incoming.TeksPesan,
			WaktuKirim: time.Now(),
			StatusBaca: false,
		}

		// Jika ada gambar terkirim (base64)
		if incoming.Base64Gambar != "" {
			folderName := fmt.Sprintf("room_%d", msg.RoomChatID)
			fileName := fmt.Sprintf("img_%d_%d", msg.UserID, time.Now().UnixNano())
			path, errSave := utils.SaveBase64Image(incoming.Base64Gambar, "chat", folderName, fileName)
			if errSave == nil {
				msg.ChatGambar = []model.ChatGambar{
					{
						PathGambar: path,
					},
				}
				// Set non-persistent field PathGambar for WebSocket broadcast JSON
				msg.PathGambar = path
			} else {
				// Log error tapi tetap lanjut simpan teks
				fmt.Printf("🔴 Gagal menyimpan gambar chat: %v\n", errSave)
			}
		}

		// Simpan pesan ke PostgreSQL lewat repository
		if err := c.repo.SaveMessage(&msg); err != nil {
			continue
		}

		// Trigger notification for the other user in the room asynchronously
		go func(msgObj model.PesanChat) {
			var room model.RoomChat
			if errRoom := config.DB.Preload("Order.Pelanggan").Preload("Order.Karyawan").First(&room, msgObj.RoomChatID).Error; errRoom == nil {
				var targetUserID uint
				var senderName string = "Seseorang"

				if msgObj.UserID == room.Order.Pelanggan.UserID {
					if room.Order.KaryawanID != nil && room.Order.Karyawan.UserID > 0 {
						targetUserID = room.Order.Karyawan.UserID
					}
					if room.Order.Pelanggan.NamaLengkap != "" {
						senderName = room.Order.Pelanggan.NamaLengkap
					}
				} else {
					targetUserID = room.Order.Pelanggan.UserID
					if room.Order.Karyawan.NamaKaryawan != "" {
						senderName = "Karyawan " + room.Order.Karyawan.NamaKaryawan
					}
				}

				shortMsg := msgObj.TeksPesan
				if shortMsg == "" {
					shortMsg = "📷 [Gambar]"
				}
				if len(shortMsg) > 60 {
					shortMsg = shortMsg[:57] + "..."
				}

				if targetUserID > 0 {
					errNotif := c.notifRepo.CreateNotificationForUser(
						targetUserID,
						fmt.Sprintf("Pesan Baru dari %s 💬", senderName),
						shortMsg,
					)
					if errNotif == nil {
						var latestNotif model.Notifikasi
						if errQuery := config.DB.Where("id_user = ?", targetUserID).Order("id_notifikasi desc").First(&latestNotif).Error; errQuery == nil {
							GlobalNotifHub.BroadcastNotification(targetUserID, latestNotif)
						}
					}
				} else {
					var adminsAndEmployees []model.User
					if errUsers := config.DB.Where("id_role IN (?)", []int{1, 2}).Find(&adminsAndEmployees).Error; errUsers == nil {
						for _, u := range adminsAndEmployees {
							errNotif := c.notifRepo.CreateNotificationForUser(
								u.IDUser,
								fmt.Sprintf("Pesan Baru dari %s 💬", senderName),
								shortMsg,
							)
							if errNotif == nil {
								var latestNotif model.Notifikasi
								if errQuery := config.DB.Where("id_user = ?", u.IDUser).Order("id_notifikasi desc").First(&latestNotif).Error; errQuery == nil {
									GlobalNotifHub.BroadcastNotification(u.IDUser, latestNotif)
								}
							}
						}
					}
				}
			}
		}(msg)

		// Kirim balik pesan ini secara broadcast ke SEMUA orang yang ada di room chat yang sama
		pool.Mu.Lock()
		for client := range pool.Clients {
			// Mengirim data pesan lengkap beserta pengirimnya ke Flutter secara real-time
			_ = client.WriteJSON(gin.H{
				"id_pesan_chat": msg.IDPesanChat,
				"id_room_chat":  msg.RoomChatID,
				"id_user":       msg.UserID,
				"teks_pesan":    msg.TeksPesan,
				"waktu_kirim":   msg.WaktuKirim.Format(time.RFC3339),
				"status_baca":   msg.StatusBaca,
				"path_gambar":   msg.PathGambar,
			})
		}
		pool.Mu.Unlock()
	}
}