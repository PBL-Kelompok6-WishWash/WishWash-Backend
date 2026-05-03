package model

import "time"

type Promo struct {
	IDPromo          uint      `gorm:"primaryKey;autoIncrement;column:id_promo"`
	KodePromo        string    `gorm:"type:varchar(50);not null;unique;column:kode_promo"`
	NamaPromo        string    `gorm:"type:varchar(150);not null;column:nama_promo"`
	Deskripsi        string    `gorm:"type:text;column:deskripsi"`
	TipePromo        string    `gorm:"type:varchar(50);column:tipe_promo"`
	NominalPotongan  float64   `gorm:"type:numeric;column:nominal_potongan"`
	MinimalOrder     float64   `gorm:"type:numeric;column:minimal_order"`
	MaksimalPotongan float64   `gorm:"type:numeric;column:maksimal_potongan"`
	TglMulai         time.Time `gorm:"type:date;column:tgl_mulai"`
	TglBerakhir      time.Time `gorm:"type:date;column:tgl_berakhir"`
	StatusPromo      string    `gorm:"type:varchar(50);column:status_promo"`
}

func (Promo) TableName() string {
	return "promo"
}