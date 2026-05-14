package model

type MetodePembayaran struct {
	IDMetodePembayaran uint   `gorm:"primaryKey;autoIncrement;column:id_metode_pembayaran" json:"id_metode_pembayaran"`
	NamaMetode         string `gorm:"type:varchar(50);not null;column:nama_metode" json:"nama_metode"`
	TipeMetode         string `gorm:"type:varchar(50);not null;column:tipe_metode" json:"tipe_metode"` // "Tunai" atau "Midtrans"
	KodeMetode         string `gorm:"type:varchar(50);column:kode_metode" json:"kode_metode"`      // Misal: "gopay", "shopeepay", "bca_va"
	GambarMetode       string `gorm:"type:text;column:gambar_metode" json:"gambar_metode"`
	StatusMetode       string `gorm:"type:varchar(20);default:'Aktif';column:status_metode" json:"status_metode"` // "Aktif" atau "Tidak Aktif"
}

func (MetodePembayaran) TableName() string {
	return "metode_pembayaran"
}
