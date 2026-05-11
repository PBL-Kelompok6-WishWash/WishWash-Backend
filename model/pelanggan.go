package model

type Pelanggan struct {
	IDPelanggan   uint   `gorm:"primaryKey;autoIncrement;column:id_pelanggan" json:"id_pelanggan"`
	UserID        uint   `gorm:"not null;column:id_user" json:"id_user"` // Nama struct diubah jadi UserID
	NamaLengkap   string `gorm:"type:varchar(150);not null;column:nama_lengkap" json:"nama_lengkap"`
	FotoPelanggan string `gorm:"type:text;column:foto_pelanggan" json:"foto_pelanggan"`
	NoTelp        string `gorm:"type:varchar(20);column:no_telp" json:"no_telp"`

	User User `gorm:"foreignKey:UserID;references:IDUser" json:"User"`
}

func (Pelanggan) TableName() string {
	return "pelanggan"
}