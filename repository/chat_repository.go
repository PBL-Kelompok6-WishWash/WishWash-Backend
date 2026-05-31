package repository

import (
	"time"

	"gorm.io/gorm" // 👈 Kita pakai gorm sekarang biar seragam!
)

// 1. Struktur mangkok data tetep sama
type MessageData struct {
	IdPesanChat int        `json:"id_pesan_chat"`
	IdRoomChat  int        `json:"id_room_chat"`
	IdUser      int        `json:"id_user"` 
	TeksPesan   string     `json:"teks_pesan"`
	WaktuKirim  time.Time  `json:"waktu_kirim"`
	StatusBaca  bool       `json:"status_baca"`
	PathGambar  *string    `json:"path_gambar"` 
}

type ChatRepository interface {
	GetMessagesByRoomID(roomID string) ([]MessageData, error)
}

type chatRepository struct {
	db *gorm.DB // 👈 Ubah dari *sql.DB ke *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository { // 👈 Ubah dari *sql.DB ke *gorm.DB
	return &chatRepository{db: db}
}

// 2. Fungsi koki menggunakan cara query GORM (Jauh lebih ringkas, gak perlu looping rows.Next!)
func (r *chatRepository) GetMessagesByRoomID(roomID string) ([]MessageData, error) {
	var messages []MessageData

	query := `
		SELECT p.id_pesan_chat, p.id_room_chat, p.id_user, p.teks_pesan, p.waktu_kirim, p.status_baca, g.path_gambar
		FROM pesan_chat p
		LEFT JOIN chat_gambar g ON p.id_pesan_chat = g.id_pesan_chat
		WHERE p.id_room_chat = ?
		ORDER BY p.waktu_kirim ASC;
	`

	// Pake r.db.Raw bawaan GORM langsung beres scan otomatis ke struct
	err := r.db.Raw(query, roomID).Scan(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}