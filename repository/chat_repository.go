package repository

import (
	"fmt"
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
	mergeDuplicateRooms(db)
	return &chatRepository{db: db}
}

func mergeDuplicateRooms(db *gorm.DB) {
	type DupeRoom struct {
		PelangganID uint
		KaryawanID  uint
		RoomCount   int
	}
	
	var dupes []DupeRoom
	db.Raw(`
		SELECT o.id_pelanggan as pelanggan_id, COALESCE(o.id_karyawan, 0) as karyawan_id, COUNT(rc.id_room_chat) as room_count
		FROM room_chat rc
		JOIN "order" o ON o.id_order = rc.id_order
		GROUP BY o.id_pelanggan, COALESCE(o.id_karyawan, 0)
		HAVING COUNT(rc.id_room_chat) > 1
	`).Scan(&dupes)

	for _, dupe := range dupes {
		var rooms []model.RoomChat
		var err error
		if dupe.KaryawanID == 0 {
			err = db.Joins("JOIN \"order\" ON \"order\".id_order = room_chat.id_order").
				Where("\"order\".id_pelanggan = ? AND \"order\".id_karyawan IS NULL", dupe.PelangganID).
				Order("room_chat.id_room_chat ASC").
				Find(&rooms).Error
		} else {
			err = db.Joins("JOIN \"order\" ON \"order\".id_order = room_chat.id_order").
				Where("\"order\".id_pelanggan = ? AND \"order\".id_karyawan = ?", dupe.PelangganID, dupe.KaryawanID).
				Order("room_chat.id_room_chat ASC").
				Find(&rooms).Error
		}

		if err == nil && len(rooms) > 1 {
			keepRoomID := rooms[0].IDRoomChat
			var deleteRoomIDs []uint
			for i := 1; i < len(rooms); i++ {
				deleteRoomIDs = append(deleteRoomIDs, rooms[i].IDRoomChat)
			}
			
			if len(deleteRoomIDs) > 0 {
				db.Model(&model.PesanChat{}).Where("id_room_chat IN ?", deleteRoomIDs).Update("id_room_chat", keepRoomID)
				db.Where("id_room_chat IN ?", deleteRoomIDs).Delete(&model.RoomChat{})
			}
		}
	}
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

	if err == nil {
		var user model.User
		r.db.First(&user, userID)

		seenPartner := make(map[string]bool)
		var uniqueRooms []model.RoomChat
		for _, room := range rooms {
			var partnerKey string
			if user.RoleID == 2 { // Karyawan, unikkan berdasarkan PelangganID
				partnerKey = fmt.Sprintf("p_%d", room.Order.PelangganID)
			} else { // Pelanggan, unikkan berdasarkan KaryawanID
				kID := uint(0)
				if room.Order.KaryawanID != nil {
					kID = *room.Order.KaryawanID
				}
				partnerKey = fmt.Sprintf("k_%d", kID)
			}

			if !seenPartner[partnerKey] {
				seenPartner[partnerKey] = true
				uniqueRooms = append(uniqueRooms, room)
			}
		}
		return uniqueRooms, nil
	}
	return rooms, err
}

// 4. Mendapatkan atau membuat Room Chat baru berdasarkan ID Order
func (r *chatRepository) GetOrCreateRoomByOrderID(orderID uint) (*model.RoomChat, error) {
	var currentOrder model.Order
	if err := r.db.Where("id_order = ?", orderID).First(&currentOrder).Error; err != nil {
		return nil, err
	}

	var existingRoom model.RoomChat
	var err error
	if currentOrder.KaryawanID == nil {
		err = r.db.Joins("JOIN \"order\" ON \"order\".id_order = room_chat.id_order").
			Where("\"order\".id_pelanggan = ? AND \"order\".id_karyawan IS NULL", currentOrder.PelangganID).
			Preload("Order").Preload("Order.Pelanggan").Preload("Order.Karyawan").
			First(&existingRoom).Error
	} else {
		err = r.db.Joins("JOIN \"order\" ON \"order\".id_order = room_chat.id_order").
			Where("\"order\".id_pelanggan = ? AND \"order\".id_karyawan = ?", currentOrder.PelangganID, *currentOrder.KaryawanID).
			Preload("Order").Preload("Order.Pelanggan").Preload("Order.Karyawan").
			First(&existingRoom).Error
	}

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