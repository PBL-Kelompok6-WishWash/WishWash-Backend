package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type AdminRepository interface {
	FindByUserID(userID uint) (*model.Admin, error)
	UpdateAdmin(admin *model.Admin) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db}
}

func (r *adminRepository) FindByUserID(userID uint) (*model.Admin, error) {
	var admin model.Admin
	// Mencari di tabel admin berdasarkan id_user
	err := r.db.Preload("User").Preload("User.Role").Where("id_user = ?", userID).First(&admin).Error
	return &admin, err
}

func (r *adminRepository) UpdateAdmin(admin *model.Admin) error {
	return r.db.Model(admin).Where("id_user = ?", admin.UserID).Update("nama_admin", admin.NamaAdmin).Error
}