package model

import "time"

type PesanChat struct {
	IDPesanChat uint      `gorm:"primaryKey;autoIncrement;column:id_pesan_chat" json:"id_pesan_chat"`
	RoomChatID  uint      `gorm:"not null;column:id_room_chat" json:"id_room_chat"`
	UserID      uint      `gorm:"not null;column:id_user" json:"id_user"` // Universal untuk semua tipe user
	TeksPesan   string    `gorm:"type:text;column:teks_pesan" json:"teks_pesan"`
	WaktuKirim  time.Time `gorm:"type:timestamp;column:waktu_kirim;default:CURRENT_TIMESTAMP" json:"waktu_kirim"`
	StatusBaca  bool      `gorm:"default:false;column:status_baca" json:"status_baca"`

	RoomChat RoomChat `gorm:"foreignKey:RoomChatID" json:"RoomChat"`
	User     User     `gorm:"foreignKey:UserID" json:"User"`
}

func (PesanChat) TableName() string {
	return "pesan_chat"
}