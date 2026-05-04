package model

import "time"

type TrackingPengiriman struct {
	IDTracking      uint      `gorm:"primaryKey;autoIncrement;column:id_tracking_pengiriman"`
	OrderID         uint      `gorm:"not null;column:id_order"`
	KaryawanID      uint      `gorm:"not null;column:id_karyawan"` // Kurir wajib ada
	LastLatitude    string    `gorm:"type:varchar(100);column:last_latitude"`
	LastLongitude   string    `gorm:"type:varchar(100);column:last_longitude"`
	DeskripsiStatus string    `gorm:"type:text;column:deskripsi_status"`
	JenisTugas      string    `gorm:"type:varchar(50);column:jenis_tugas"` // Contoh: "pickup" atau "delivery"
	WaktuUpdate     time.Time `gorm:"type:timestamp;column:waktu_update;default:CURRENT_TIMESTAMP"`

	Order    Order    `gorm:"foreignKey:OrderID"`
	Karyawan Karyawan `gorm:"foreignKey:KaryawanID"`
}

func (TrackingPengiriman) TableName() string {
	return "tracking_pengiriman"
}