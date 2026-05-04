package model

import "time"

type PesanChat struct {
	IDPesanChat uint      `gorm:"primaryKey;autoIncrement;column:id_pesan_chat"`
	RoomChatID  uint      `gorm:"not null;column:id_room_chat"`
	UserID      uint      `gorm:"not null;column:id_user"` // Universal untuk semua tipe user
	TeksPesan   string    `gorm:"type:text;column:teks_pesan"`
	WaktuKirim  time.Time `gorm:"type:timestamp;column:waktu_kirim;default:CURRENT_TIMESTAMP"`
	StatusBaca  bool      `gorm:"default:false;column:status_baca"`

	RoomChat RoomChat `gorm:"foreignKey:RoomChatID"`
	User     User     `gorm:"foreignKey:UserID"`
}

func (PesanChat) TableName() string {
	return "pesan_chat"
}