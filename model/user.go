package model

type User struct {
	IDUser   uint   `gorm:"primaryKey;autoIncrement;column:id_user"`
	RoleID   uint   `gorm:"not null;column:id_role"` // Nama struct diubah jadi RoleID agar GORM paham
	Username string `gorm:"type:varchar(100);not null;unique;column:username"`
	Email    string `gorm:"type:varchar(100);not null;unique;column:email"`
	Password string `gorm:"type:varchar(255);not null;column:password"`

	// Relasi ke tabel Role
	Role Role `gorm:"foreignKey:RoleID;references:IDRole"`
}

func (User) TableName() string {
	return "user"
}