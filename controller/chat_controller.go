package controller

import (
	"net/http"

	"github.com/PBL-Kelompok6-WishWash/backend/repository"
	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatRepo repository.ChatRepository
}

// Constructor buat bikin instance controller baru
func NewChatController(chatRepo repository.ChatRepository) *ChatController {
	return &ChatController{chatRepo: chatRepo}
}

// Fungsi utama untuk handle request GET riwayat pesan
func (h *ChatController) GetMessages(c *gin.Context) {
	// 1. Tangkap parameter id_room_chat dari URL rute main.go
	roomID := c.Param("id_room_chat")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID room chat tidak valid atau kosong",
		})
		return
	}

	// 2. Minta si koki database (repository) buat narik data dari PostgreSQL
	messages, err := h.chatRepo.GetMessagesByRoomID(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat pesan dari database",
			"error":   err.Error(),
		})
		return
	}

	// 3. Kalau datanya kosong (belum pernah chat), jangan eror, tapi kasih array kosong []
	if messages == nil {
		messages = []repository.MessageData{}
	}

	// 4. Sajikan data chat-nya dalam bentuk JSON yang super rapi!
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil memuat riwayat pesan",
		"data":    messages,
	})
}