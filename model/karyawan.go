package model

type Karyawan struct {
	IDKaryawan         uint   `gorm:"primaryKey;autoIncrement;column:id_karyawan" json:"id_karyawan"`
	UserID             uint   `gorm:"not null;column:id_user" json:"id_user"`
	NamaKaryawan       string `gorm:"type:varchar(150);not null;column:nama_karyawan" json:"nama_karyawan"`
	FotoKaryawan       string `gorm:"type:text;column:foto_karyawan" json:"foto_karyawan"`
	NoTelp             string `gorm:"type:varchar(20);column:no_telp" json:"no_telp"`
	PlatNomor          string `gorm:"type:varchar(20);column:plat_nomor" json:"plat_nomor"`
	JenisKendaraan     string `gorm:"type:varchar(50);column:jenis_kendaraan" json:"jenis_kendaraan"`
	StatusKetersediaan string `gorm:"type:varchar(50);column:status_ketersediaan" json:"status_ketersediaan"`

	User User `gorm:"foreignKey:UserID;references:IDUser" json:"User"`
}

func (Karyawan) TableName() string {
	return "karyawan"
}