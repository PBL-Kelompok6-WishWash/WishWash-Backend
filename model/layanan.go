package model

type Layanan struct {
	IDLayanan      uint    `gorm:"primaryKey;autoIncrement;column:id_layanan"`
	NamaLayanan    string  `gorm:"type:varchar(100);not null;column:nama_layanan"`
	GambarLayanan  string  `gorm:"type:text;column:gambar_layanan"`
	JenisSatuan    string  `gorm:"type:varchar(50);column:jenis_satuan"`
	HargaPerSatuan float64 `gorm:"type:numeric;column:harga_per_satuan"`
}

func (Layanan) TableName() string {
	return "layanan"
}