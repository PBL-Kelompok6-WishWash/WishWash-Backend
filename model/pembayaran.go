package model

import "time"

type Pembayaran struct {
	IDPembayaran     uint      `gorm:"primaryKey;autoIncrement;column:id_pembayaran"`
	OrderID          uint      `gorm:"not null;column:id_order"`
	KaryawanID       *uint     `gorm:"column:id_karyawan"` // Boleh null jika bayar online (sistem)
	MetodeBayar      string    `gorm:"type:varchar(50);column:metode_bayar"`
	JumlahBayar      float64   `gorm:"type:numeric;not null;column:jumlah_bayar"`
	StatusPembayaran string    `gorm:"type:varchar(50);column:status_pembayaran"`
	ReferensiBayar   string    `gorm:"type:varchar(100);column:referensi_bayar"`
	TglPembayaran    time.Time `gorm:"type:timestamp;column:tgl_pembayaran"`

	Order    Order    `gorm:"foreignKey:OrderID"`
	Karyawan Karyawan `gorm:"foreignKey:KaryawanID"`
}

func (Pembayaran) TableName() string {
	return "pembayaran"
}