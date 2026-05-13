package model

type Layanan struct {
	IDLayanan      uint    `gorm:"primaryKey;autoIncrement;column:id_layanan" json:"id_layanan"`
	NamaLayanan    string  `gorm:"type:varchar(100);not null;column:nama_layanan" json:"nama_layanan"`
	GambarLayanan  string  `gorm:"type:text;column:gambar_layanan" json:"gambar_layanan"`
	JenisSatuan    string  `gorm:"type:varchar(50);column:jenis_satuan" json:"jenis_satuan"`
	HargaPerSatuan float64 `gorm:"type:numeric;column:harga_per_satuan" json:"harga_per_satuan"`

	// Relasi ke ReferensiStatusLayanan (1 Layanan -> N Status)
	ReferensiStatus []ReferensiStatusLayanan `gorm:"foreignKey:LayananID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"referensi_status"`

	// Relasi ke PaketLayanan (1 Layanan -> N PaketLayanan)
	PaketLayanan []PaketLayanan `gorm:"foreignKey:LayananID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"paket_layanan"`
}

func (Layanan) TableName() string {
	return "layanan"
}