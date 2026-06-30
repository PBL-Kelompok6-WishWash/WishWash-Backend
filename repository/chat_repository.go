package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	GetMessagesByRoomID(roomID uint) ([]model.PesanChat, error)
	SaveMessage(msg *model.PesanChat) error
	GetRoomsByUserID(userID uint) ([]model.RoomChat, error) // Tambahan untuk daftar chat di awal
	GetOrCreateRoomByOrderID(orderID uint) (*model.RoomChat, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func mergeDuplicateRooms(db *gorm.DB) {
	// Fungsi ini dinonaktifkan agar chat room per order / per karyawan tidak digabungkan
}

// 1. Mengambil riwayat pesan lama berdasarkan ID Room Chat
func (r *chatRepository) GetMessagesByRoomID(roomID uint) ([]model.PesanChat, error) {
	var messages []model.PesanChat
	// Preload "User" dan "ChatGambar" untuk data lengkap
	err := r.db.Where("id_room_chat = ?", roomID).Order("waktu_kirim asc").Preload("User").Preload("ChatGambar").Find(&messages).Error
	if err == nil {
		for i := range messages {
			if len(messages[i].ChatGambar) > 0 {
				messages[i].PathGambar = messages[i].ChatGambar[0].PathGambar
			}
		}
	}
	return messages, nil
}

// 2. Menyimpan pesan baru yang masuk lewat WebSocket ke database
func (r *chatRepository) SaveMessage(msg *model.PesanChat) error {
	return r.db.Create(msg).Error
}

// 3. Mengambil daftar Room Chat yang aktif untuk user tertentu (Pelanggan/Karyawan)
func (r *chatRepository) GetRoomsByUserID(userID uint) ([]model.RoomChat, error) {
	var rooms []model.RoomChat
	// Mencari room chat yang terikat dengan order milik pelanggan atau ditangani karyawan tersebut
	err := r.db.Joins("JOIN \"order\" ON \"order\".id_order = room_chat.id_order").
		Joins("JOIN pelanggan ON pelanggan.id_pelanggan = \"order\".id_pelanggan").
		Where("pelanggan.id_user = ? OR \"order\".id_karyawan = (SELECT id_karyawan FROM karyawan WHERE id_user = ?)", userID, userID).
		Preload("Order").Preload("Order.Pelanggan").Preload("Order.Karyawan").
		Order("room_chat.waktu_dibuat DESC").
		Find(&rooms).Error

	return rooms, err
}

// 4. Mendapatkan atau membuat Room Chat baru berdasarkan ID Order
func (r *chatRepository) GetOrCreateRoomByOrderID(orderID uint) (*model.RoomChat, error) {
	var existingRoom model.RoomChat
	err := r.db.Where("id_order = ?", orderID).
		Preload("Order").Preload("Order.Pelanggan").Preload("Order.Karyawan").
		First(&existingRoom).Error

	if err == nil {
		return &existingRoom, nil
	}

	room := model.RoomChat{
		OrderID: orderID,
	}
	if err := r.db.Create(&room).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Order").Preload("Order.Pelanggan").Preload("Order.Karyawan").First(&room)
	return &room, nil
}