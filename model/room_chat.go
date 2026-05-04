package model

import "time"

type RoomChat struct {
	IDRoomChat  uint      `gorm:"primaryKey;autoIncrement;column:id_room_chat"`
	OrderID     uint      `gorm:"not null;column:id_order"` // Typo di ERD (id_oder) sudah diluruskan
	WaktuDibuat time.Time `gorm:"type:timestamp;column:waktu_dibuat;default:CURRENT_TIMESTAMP"`

	Order Order `gorm:"foreignKey:OrderID"`
}

func (RoomChat) TableName() string {
	return "room_chat"
}