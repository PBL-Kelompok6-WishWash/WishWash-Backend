package model

type PaketLayanan struct {
	IDPaketLayanan uint    `gorm:"primaryKey;autoIncrement;column:id_paket_layanan" json:"id_paket_layanan"`
	NamaPaket      string  `gorm:"type:varchar(100);not null;column:nama_paket" json:"nama_paket"`
	DurasiJam      int     `gorm:"column:durasi_jam" json:"durasi_jam"`
	BiayaTambahan  float64 `gorm:"type:numeric;column:biaya_tambahan" json:"biaya_tambahan"`
}

func (PaketLayanan) TableName() string {
	return "paket_layanan"
}