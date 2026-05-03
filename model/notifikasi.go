package model

import "time"

type Notifikasi struct {
	IDNotifikasi uint      `gorm:"primaryKey;autoIncrement;column:id_notifikasi"`
	UserID       uint      `gorm:"not null;column:id_user"` // FK ke tabel user
	Judul        string    `gorm:"type:varchar(150);not null;column:judul"`
	Pesan        string    `gorm:"type:text;not null;column:pesan"`
	IsRead       bool      `gorm:"default:false;column:is_read"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`

	User User `gorm:"foreignKey:UserID"`
}

func (Notifikasi) TableName() string {
	return "notifikasi"
}