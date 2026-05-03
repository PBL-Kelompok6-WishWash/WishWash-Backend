package model

type Parfum struct {
	IDParfum   uint   `gorm:"primaryKey;autoIncrement;column:id_parfum"`
	NamaParfum string `gorm:"type:varchar(100);not null;column:nama_parfum"`
	Keterangan string `gorm:"type:text;column:keterangan"` 
}

func (Parfum) TableName() string {
	return "parfum"
}