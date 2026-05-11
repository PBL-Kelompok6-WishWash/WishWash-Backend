package model

type User struct {
	IDUser   uint   `gorm:"primaryKey;autoIncrement;column:id_user" json:"id_user"`
	RoleID   uint   `gorm:"not null;column:id_role" json:"id_role"`
	Username string `gorm:"type:varchar(100);not null;unique;column:username" json:"username"`
	Email    string `gorm:"type:varchar(100);not null;unique;column:email" json:"email"`
	Password string `gorm:"type:varchar(255);not null;column:password" json:"password"`

	// Relasi ke tabel Role
	Role Role `gorm:"foreignKey:RoleID;references:IDRole" json:"Role"`
}

func (User) TableName() string {
	return "user"
}