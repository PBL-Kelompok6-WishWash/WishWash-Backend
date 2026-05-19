package model

type Alamat struct {
	IDAlamat      uint    `gorm:"primaryKey;autoIncrement;column:id_alamat" json:"id_alamat"`
	PelangganID   uint    `gorm:"not null;column:id_pelanggan" json:"id_pelanggan"` // FK ke tabel pelanggan (menggantikan id_customer di ERD)
	Latitude      string  `gorm:"type:varchar(100);column:latitude" json:"latitude"` // Bisa pakai float, tapi varchar lebih aman untuk presisi map API
	Longitude     string  `gorm:"type:varchar(100);column:longitude" json:"longitude"`
	AlamatLengkap string  `gorm:"type:text;not null;column:alamat_lengkap" json:"alamat_lengkap"`
	TipeAlamat    string  `gorm:"type:varchar(50);column:tipe_alamat" json:"tipe_alamat"` // Misal: Rumah, Kos, Kantor
	NamaPenerima  string  `gorm:"type:varchar(100);column:nama_penerima" json:"nama_penerima"`
	NoHpPenerima  string  `gorm:"type:varchar(20);column:nohp_penerima" json:"nohp_penerima"`
	IsPrimary     bool    `gorm:"type:boolean;default:false;column:is_primary" json:"is_primary"`

	Pelanggan *Pelanggan `gorm:"foreignKey:PelangganID" json:"Pelanggan,omitempty"`
}

func (Alamat) TableName() string {
	return "alamat"
}