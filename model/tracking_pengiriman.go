package model

import "time"

type TrackingPengiriman struct {
	IDTracking      uint      `gorm:"primaryKey;autoIncrement;column:id_tracking_pengiriman" json:"id_tracking_pengiriman"`
	OrderID         uint      `gorm:"not null;column:id_order" json:"id_order"`
	KaryawanID      uint      `gorm:"not null;column:id_karyawan" json:"id_karyawan"` // Kurir wajib ada
	LastLatitude    string    `gorm:"type:varchar(100);column:last_latitude" json:"last_latitude"`
	LastLongitude   string    `gorm:"type:varchar(100);column:last_longitude" json:"last_longitude"`
	DeskripsiStatus string    `gorm:"type:text;column:deskripsi_status" json:"deskripsi_status"`
	JenisTugas      string    `gorm:"type:varchar(50);column:jenis_tugas" json:"jenis_tugas"` // Contoh: "pickup" atau "delivery"
	WaktuUpdate     time.Time `gorm:"type:timestamp;column:waktu_update;default:CURRENT_TIMESTAMP" json:"waktu_update"`

	Order    Order    `gorm:"foreignKey:OrderID" json:"Order"`
	Karyawan Karyawan `gorm:"foreignKey:KaryawanID" json:"Karyawan"`
}

func (TrackingPengiriman) TableName() string {
	return "tracking_pengiriman"
}