package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type KaryawanRepository interface {
	CreateKaryawan(karyawan *model.Karyawan) error
	FindByUserID(userID uint) (*model.Karyawan, error)
	UpdateKaryawan(karyawan *model.Karyawan) error
}

type karyawanRepository struct {
	db *gorm.DB
}

func NewKaryawanRepository(db *gorm.DB) KaryawanRepository {
	return &karyawanRepository{db}
}

func (r *karyawanRepository) CreateKaryawan(karyawan *model.Karyawan) error {
	return r.db.Create(karyawan).Error
}

func (r *karyawanRepository) FindByUserID(userID uint) (*model.Karyawan, error) {
	var karyawan model.Karyawan
	err := r.db.Where("id_user = ?", userID).First(&karyawan).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *karyawanRepository) UpdateKaryawan(karyawan *model.Karyawan) error {
	return r.db.Model(karyawan).Where("id_user = ?", karyawan.UserID).Update("nama_karyawan", karyawan.NamaKaryawan).Error
}
