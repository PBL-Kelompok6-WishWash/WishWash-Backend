package model

type ChatGambar struct {
	IDChatGambar uint   `gorm:"primaryKey;autoIncrement;column:id_chat_gambar"`
	PesanChatID  uint   `gorm:"not null;column:id_pesan_chat"`
	PathGambar   string `gorm:"type:text;not null;column:path_gambar"`

	PesanChat PesanChat `gorm:"foreignKey:PesanChatID"`
}

func (ChatGambar) TableName() string {
	return "chat_gambar"
} 