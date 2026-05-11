package model

import "time"

type RiwayatStatusDetail struct {
	IDRiwayat         uint      `gorm:"primaryKey;autoIncrement;column:id_riwayat_status_detail" json:"id_riwayat_status_detail"`
	ReferensiStatusID uint      `gorm:"not null;column:id_referensi_status_layanan" json:"id_referensi_status_layanan"`
	OrderID           uint      `gorm:"not null;column:id_order" json:"id_order"`
	KaryawanID        *uint     `gorm:"column:id_karyawan" json:"id_karyawan"` // Pointer, karena kadang sistem otomatis update status tanpa karyawan
	WaktuUpdate       time.Time `gorm:"type:timestamp;column:waktu_update;default:CURRENT_TIMESTAMP" json:"waktu_update"`

	ReferensiStatus ReferensiStatusLayanan `gorm:"foreignKey:ReferensiStatusID" json:"ReferensiStatus"`
	Order           Order                  `gorm:"foreignKey:OrderID" json:"Order"`
	Karyawan        Karyawan               `gorm:"foreignKey:KaryawanID" json:"Karyawan"`
}

func (RiwayatStatusDetail) TableName() string {
	return "riwayat_status_detail"
}