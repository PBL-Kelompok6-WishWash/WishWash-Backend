package model

type ReferensiStatusLayanan struct {
	IDReferensiStatus uint   `gorm:"primaryKey;autoIncrement;column:id_referensi_status_layanan"`
	LayananID         uint   `gorm:"not null;column:id_layanan"` // FK ke tabel layanan
	NamaStatus        string `gorm:"type:varchar(100);not null;column:nama_status"`
	UrutanTahap       int    `gorm:"not null;column:urutan_tahap"`

	Layanan Layanan `gorm:"foreignKey:LayananID"`
}

func (ReferensiStatusLayanan) TableName() string {
	return "referensi_status_layanan"
}