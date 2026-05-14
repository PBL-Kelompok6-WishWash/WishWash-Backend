package model

type Parfum struct {
	IDParfum   uint   `gorm:"primaryKey;autoIncrement;column:id_parfum" json:"id_parfum"`
	NamaParfum string `gorm:"type:varchar(100);not null;column:nama_parfum" json:"nama_parfum"`
	Keterangan   string `gorm:"type:text;column:keterangan" json:"keterangan"`
	StatusParfum string `gorm:"type:varchar(20);not null;default:'Tersedia';column:status_parfum" json:"status_parfum"`
}

func (Parfum) TableName() string {
	return "parfum"
}