package repository

import (
	"log"

	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type NotifikasiRepository interface {
	FindAllByUserID(userID uint) ([]model.Notifikasi, error)
	MarkAsRead(id uint, userID uint) error
	MarkAllAsRead(userID uint) error
	CreateNotificationForAdmins(title string, message string) error
}

type notifikasiRepository struct {
	db *gorm.DB
}

func NewNotifikasiRepository(db *gorm.DB) NotifikasiRepository {
	return &notifikasiRepository{db}
}

func (r *notifikasiRepository) FindAllByUserID(userID uint) ([]model.Notifikasi, error) {
	var notifications []model.Notifikasi
	err := r.db.Where("id_user = ?", userID).Order("id_notifikasi DESC").Find(&notifications).Error
	return notifications, err
}

func (r *notifikasiRepository) MarkAsRead(id uint, userID uint) error {
	return r.db.Model(&model.Notifikasi{}).
		Where("id_notifikasi = ? AND id_user = ?", id, userID).
		Update("is_read", true).Error
}

func (r *notifikasiRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&model.Notifikasi{}).
		Where("id_user = ?", userID).
		Update("is_read", true).Error
}

func (r *notifikasiRepository) CreateNotificationForAdmins(title string, message string) error {
	var admins []model.User
	if err := r.db.Where("id_role = ?", 1).Find(&admins).Error; err != nil {
		return err
	}

	for _, admin := range admins {
		notif := model.Notifikasi{
			UserID: admin.IDUser,
			Judul:  title,
			Pesan:  message,
			IsRead: false,
		}
		if err := r.db.Create(&notif).Error; err != nil {
			log.Printf("⚠️ Gagal membuat notifikasi untuk admin UserID %d: %v", admin.IDUser, err)
		}
	}
	return nil
}
