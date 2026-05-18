package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type KaryawanRepository interface {
	CreateKaryawan(karyawan *model.Karyawan) error
	FindByUserID(userID uint) (*model.Karyawan, error)
	UpdateKaryawan(karyawan *model.Karyawan) error

	FindAll() ([]model.Karyawan, error)
	FindByID(idKaryawan uint) (*model.Karyawan, error)
	Update(karyawan *model.Karyawan) error
	Delete(idKaryawan uint) error
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
	err := r.db.Preload("User").Preload("User.Role").Where("id_user = ?", userID).First(&karyawan).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *karyawanRepository) UpdateKaryawan(karyawan *model.Karyawan) error {
	return r.db.Model(&model.Karyawan{}).
		Where("id_user = ?", karyawan.UserID).
		Updates(map[string]interface{}{
			"nama_karyawan":   karyawan.NamaKaryawan,
			"no_telp":         karyawan.NoTelp,
			"foto_karyawan":   karyawan.FotoKaryawan,
			"plat_nomor":      karyawan.PlatNomor,
			"jenis_kendaraan": karyawan.JenisKendaraan,
		}).Error
}

func (r *karyawanRepository) FindAll() ([]model.Karyawan, error) {
	var karyawans []model.Karyawan
	err := r.db.Preload("User").Preload("User.Role").Find(&karyawans).Error
	return karyawans, err
}

func (r *karyawanRepository) FindByID(idKaryawan uint) (*model.Karyawan, error) {
	var karyawan model.Karyawan
	err := r.db.Preload("User").Preload("User.Role").First(&karyawan, idKaryawan).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *karyawanRepository) Update(karyawan *model.Karyawan) error {
	return r.db.Save(karyawan).Error
}

func (r *karyawanRepository) Delete(idKaryawan uint) error {
	return r.db.Delete(&model.Karyawan{}, idKaryawan).Error
}
