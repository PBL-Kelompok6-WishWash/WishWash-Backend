package model

import "time"

type Promo struct {
	IDPromo          uint      `gorm:"primaryKey;autoIncrement;column:id_promo" json:"id_promo"`
	KodePromo        string    `gorm:"type:varchar(50);not null;unique;column:kode_promo" json:"kode_promo"`
	NamaPromo        string    `gorm:"type:varchar(150);not null;column:nama_promo" json:"nama_promo"`
	Deskripsi        string    `gorm:"type:text;column:deskripsi" json:"deskripsi"`
	TipePromo        string    `gorm:"type:varchar(50);column:tipe_promo" json:"tipe_promo"`
	NominalPotongan  float64   `gorm:"type:numeric;column:nominal_potongan" json:"nominal_potongan"`
	MinimalOrder     float64   `gorm:"type:numeric;column:minimal_order" json:"minimal_order"`
	MaksimalPotongan float64   `gorm:"type:numeric;column:maksimal_potongan" json:"maksimal_potongan"`
	TglMulai         time.Time `gorm:"type:date;column:tgl_mulai" json:"tgl_mulai"`
	TglBerakhir      time.Time `gorm:"type:date;column:tgl_berakhir" json:"tgl_berakhir"`
	StatusPromo      string    `gorm:"type:varchar(50);column:status_promo" json:"status_promo"`
	GambarPromo      string    `gorm:"type:text;column:gambar_promo" json:"gambar_promo"`
}

func (Promo) TableName() string {
	return "promo"
}