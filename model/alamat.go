package model

type Alamat struct {
	IDAlamat      uint    `gorm:"primaryKey;autoIncrement;column:id_alamat"`
	PelangganID   uint    `gorm:"not null;column:id_pelanggan"` // FK ke tabel pelanggan (menggantikan id_customer di ERD)
	Latitude      string  `gorm:"type:varchar(100);column:latitude"` // Bisa pakai float, tapi varchar lebih aman untuk presisi map API
	Longitude     string  `gorm:"type:varchar(100);column:longitude"`
	AlamatLengkap string  `gorm:"type:text;not null;column:alamat_lengkap"`
	TipeAlamat    string  `gorm:"type:varchar(50);column:tipe_alamat"` // Misal: Rumah, Kos, Kantor

	Pelanggan Pelanggan `gorm:"foreignKey:PelangganID"`
}

func (Alamat) TableName() string {
	return "alamat"
}