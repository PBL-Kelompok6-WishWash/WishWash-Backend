package model

import "time"

type RiwayatStatusDetail struct {
	IDRiwayat         uint      `gorm:"primaryKey;autoIncrement;column:id_riwayat_status_detail"`
	ReferensiStatusID uint      `gorm:"not null;column:id_referensi_status_layanan"`
	OrderID           uint      `gorm:"not null;column:id_order"`
	KaryawanID        *uint     `gorm:"column:id_karyawan"` // Pointer, karena kadang sistem otomatis update status tanpa karyawan
	WaktuUpdate       time.Time `gorm:"type:timestamp;column:waktu_update;default:CURRENT_TIMESTAMP"`

	ReferensiStatus ReferensiStatusLayanan `gorm:"foreignKey:ReferensiStatusID"`
	Order           Order                  `gorm:"foreignKey:OrderID"`
	Karyawan        Karyawan               `gorm:"foreignKey:KaryawanID"`
}

func (RiwayatStatusDetail) TableName() string {
	return "riwayat_status_detail"
}