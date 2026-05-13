package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

// 1. Interface: Kontrak kerja untuk Pelanggan (Lama + Baru)
type PelangganRepository interface {
	CreatePelanggan(pelanggan *model.Pelanggan) error
	FindByUserID(userID uint) (*model.Pelanggan, error)
	UpdatePelanggan(pelanggan *model.Pelanggan) error

	FindAll() ([]model.Pelanggan, error)
	FindByID(idPelanggan uint) (*model.Pelanggan, error)
	Update(pelanggan *model.Pelanggan) error
	Delete(idPelanggan uint) error
}

// 2. Struct
type pelangganRepository struct {
	db *gorm.DB
}

// 3. Constructor
func NewPelangganRepository(db *gorm.DB) PelangganRepository {
	return &pelangganRepository{db}
}

// =========================================================
// IMPLEMENTASI FUNGSI LAMA (Auth & Profile)
// =========================================================

func (r *pelangganRepository) CreatePelanggan(pelanggan *model.Pelanggan) error {
	return r.db.Create(pelanggan).Error
}

func (r *pelangganRepository) FindByUserID(userID uint) (*model.Pelanggan, error) {
	var pelanggan model.Pelanggan
	err := r.db.Preload("User").Preload("User.Role").Where("id_user = ?", userID).First(&pelanggan).Error
	if err != nil {
		return nil, err
	}
	return &pelanggan, nil
}

func (r *pelangganRepository) UpdatePelanggan(pelanggan *model.Pelanggan) error {
	return r.db.Model(pelanggan).Where("id_user = ?", pelanggan.UserID).Update("nama_lengkap", pelanggan.NamaLengkap).Error
}

func (r *pelangganRepository) FindAll() ([]model.Pelanggan, error) {
	var pelanggans []model.Pelanggan
	// Preload("User") agar username & email ikut terambil otomatis
	err := r.db.Preload("User").Find(&pelanggans).Error
	return pelanggans, err
}

func (r *pelangganRepository) FindByID(idPelanggan uint) (*model.Pelanggan, error) {
	var pelanggan model.Pelanggan
	err := r.db.Preload("User").First(&pelanggan, idPelanggan).Error
	if err != nil {
		return nil, err
	}
	return &pelanggan, nil
}

func (r *pelangganRepository) Update(pelanggan *model.Pelanggan) error {
	// Save() akan meng-update seluruh field di database sesuai isi struct terbaru
	return r.db.Save(pelanggan).Error
}

func (r *pelangganRepository) Delete(idPelanggan uint) error {
	// Menghapus data pelanggan berdasarkan ID Pelanggan
	return r.db.Delete(&model.Pelanggan{}, idPelanggan).Error
}