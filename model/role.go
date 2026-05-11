package model

// Struct Role merepresentasikan tabel 'roles' di database
type Role struct {
	IDRole   uint   `gorm:"primaryKey;autoIncrement;column:id_role" json:"id_role"`
	NamaRole string `gorm:"type:varchar(50);not null;column:nama_role" json:"nama_role"`
}

// Memaksa GORM agar menggunakan nama tabel 'roles' (sesuai ERD)
func (Role) TableName() string {
	return "roles"
}