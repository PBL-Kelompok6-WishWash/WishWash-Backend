package model

import "time"

type RoomChat struct {
	IDRoomChat  uint      `gorm:"primaryKey;autoIncrement;column:id_room_chat" json:"id_room_chat"`
	OrderID     uint      `gorm:"not null;column:id_order" json:"id_order"` // Typo di ERD (id_oder) sudah diluruskan
	WaktuDibuat time.Time `gorm:"type:timestamp;column:waktu_dibuat;default:CURRENT_TIMESTAMP" json:"waktu_dibuat"`

	Order Order `gorm:"foreignKey:OrderID" json:"Order"`
}

func (RoomChat) TableName() string {
	return "room_chat"
}