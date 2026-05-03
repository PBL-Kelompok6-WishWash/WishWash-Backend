package model

type Pelanggan struct {
	IDPelanggan   uint   `gorm:"primaryKey;autoIncrement;column:id_pelanggan"`
	UserID        uint   `gorm:"not null;column:id_user"` // Nama struct diubah jadi UserID
	NamaLengkap   string `gorm:"type:varchar(150);not null;column:nama_lengkap"`
	FotoPelanggan string `gorm:"type:text;column:foto_pelanggan"`
	NoTelp        string `gorm:"type:varchar(20);column:no_telp"`

	User User `gorm:"foreignKey:UserID;references:IDUser"`
}

func (Pelanggan) TableName() string {
	return "pelanggan"
}