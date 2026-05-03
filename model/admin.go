package model

type Admin struct {
	IDAdmin   uint   `gorm:"primaryKey;autoIncrement;column:id_admin"`
	UserID    uint   `gorm:"not null;column:id_user"` // Nama struct diubah jadi UserID
	NamaAdmin string `gorm:"type:varchar(100);not null;column:nama_admin"`

	User User `gorm:"foreignKey:UserID;references:IDUser"`
}

func (Admin) TableName() string {
	return "admin"
}