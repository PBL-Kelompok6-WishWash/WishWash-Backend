package model

import "time"

type Penilaian struct {
	IDPenilaian  uint      `gorm:"primaryKey;autoIncrement;column:id_penilaian" json:"id_penilaian"`
	OrderID      uint      `gorm:"not null;column:id_order" json:"id_order"`
	Ulasan       string    `gorm:"type:text;column:ulasan" json:"ulasan"`
	Bintang          int       `gorm:"not null;column:bintang" json:"bintang"`
	BintangLayanan   int       `gorm:"not null;column:bintang_layanan;default:5" json:"bintang_layanan"`
	BintangKurir     int       `gorm:"not null;column:bintang_kurir;default:5" json:"bintang_kurir"`
	BintangKecepatan int       `gorm:"not null;column:bintang_kecepatan;default:5" json:"bintang_kecepatan"`
	TglPenilaian time.Time `gorm:"type:timestamp;column:tgl_penilaian;default:CURRENT_TIMESTAMP" json:"tgl_penilaian"`

	Order Order `gorm:"foreignKey:OrderID" json:"Order"`
}

func (Penilaian) TableName() string {
	return "penilaian"
}