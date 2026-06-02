package controller

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 1. Konfigurasi Upgrader buat ngubah HTTP biasa jadi WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Biar Flutter bisa konek tanpa kehalang CORS
	},
}

// 2. Map buat nyatet siapa aja yang lagi online di room chat tertentu
// Struktur: map[id_room_chat][]koneksi_websocket
var (
	clients   = make(map[string][]*websocket.Conn)
	clientsMu sync.Mutex // Biar gak tabrakan memori saat banyak yang chat bareng
)

// HandleWS adalah fungsi utama pelayan WebSocket kita
func (h *ChatController) HandleWS(c *gin.Context) {
	// Ambil ID Room dari URL
	roomID := c.Param("id_room_chat")
	if roomID == "" {
		log.Println("🔴 WS Error: Room ID kosong")
		return
	}

	// 3. Upgrade koneksi HTTP ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("🔴 Gagal upgrade ke Websocket: %v\n", err)
		return
	}
	defer conn.Close()

	// 4. Masukkan koneksi user ini ke dalam daftar room yang aktif
	clientsMu.Lock()
	clients[roomID] = append(clients[roomID], conn)
	clientsMu.Unlock()

	log.Printf("📱 User terhubung ke WebSocket Room #%s\n", roomID)

	// 5. Loop ini bakal jalan terus selama user masih buka halaman chat
	for {
		var msg map[string]interface{}

		// Dengerin apakah ada pesan JSON masuk dari Flutter
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("🔌 User keluar/terputus dari room %s\n", roomID)
			break
		}

		log.Printf("📩 Pesan Baru di Room %s: %v\n", roomID, msg)

		// 6. BROADCAST: Kirim pesan ini ke semua orang yang lagi buka room yang sama
		clientsMu.Lock()
		for _, client := range clients[roomID] {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("🔴 Gagal kirim broadcast ke salah satu user: %v\n", err)
				client.Close()
			}
		}
		clientsMu.Unlock()
	}

	// 7. Kalau user nutup aplikasi / keluar chat, hapus koneksinya dari daftar biar gak bocor memorinya
	clientsMu.Lock()
	for i, client := range clients[roomID] {
		if client == conn {
			clients[roomID] = append(clients[roomID][:i], clients[roomID][i+1:]...)
			break
		}
	}
	clientsMu.Unlock()
}