package model

import "time"

type Pembayaran struct {
	IDPembayaran     uint      `gorm:"primaryKey;autoIncrement;column:id_pembayaran" json:"id_pembayaran"`
	OrderID          uint      `gorm:"not null;column:id_order" json:"id_order"`
	KaryawanID       *uint     `gorm:"column:id_karyawan" json:"id_karyawan"` // Boleh null jika bayar online (sistem)
	MetodeBayar      string    `gorm:"type:varchar(50);column:metode_bayar" json:"metode_bayar"`
	JumlahBayar      float64   `gorm:"type:numeric;not null;column:jumlah_bayar" json:"jumlah_bayar"`
	StatusPembayaran string    `gorm:"type:varchar(50);column:status_pembayaran" json:"status_pembayaran"`
	ReferensiBayar   string    `gorm:"type:varchar(100);column:referensi_bayar" json:"referensi_bayar"`
	TglPembayaran    time.Time `gorm:"type:timestamp;column:tgl_pembayaran" json:"tgl_pembayaran"`

	Order    Order    `gorm:"foreignKey:OrderID" json:"Order"`
	Karyawan Karyawan `gorm:"foreignKey:KaryawanID" json:"Karyawan"`
}

func (Pembayaran) TableName() string {
	return "pembayaran"
}