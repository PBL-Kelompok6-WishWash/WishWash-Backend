package model

import "time"

type Penilaian struct {
	IDPenilaian  uint      `gorm:"primaryKey;autoIncrement;column:id_penilaian"`
	OrderID      uint      `gorm:"not null;column:id_order"`
	Ulasan       string    `gorm:"type:text;column:ulasan"`
	Bintang      int       `gorm:"not null;column:bintang"`
	TglPenilaian time.Time `gorm:"type:timestamp;column:tgl_penilaian;default:CURRENT_TIMESTAMP"`

	Order Order `gorm:"foreignKey:OrderID"`
}

func (Penilaian) TableName() string {
	return "penilaian"
}