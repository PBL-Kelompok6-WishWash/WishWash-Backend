package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

// 1. Interface: Kontrak kerja untuk Pelanggan
type PelangganRepository interface {
	CreatePelanggan(pelanggan *model.Pelanggan) error
	FindByUserID(userID uint) (*model.Pelanggan, error)
}

// 2. Struct
type pelangganRepository struct {
	db *gorm.DB
}

// 3. Constructor
func NewPelangganRepository(db *gorm.DB) PelangganRepository {
	return &pelangganRepository{db}
}

// 4. Implementasi CreatePelanggan: Dipakai saat Register
func (r *pelangganRepository) CreatePelanggan(pelanggan *model.Pelanggan) error {
	return r.db.Create(pelanggan).Error
}

// 5. Implementasi FindByUserID: Mencari profil pelanggan berdasarkan ID User-nya
func (r *pelangganRepository) FindByUserID(userID uint) (*model.Pelanggan, error) {
	var pelanggan model.Pelanggan

	// Cari di database: "SELECT * FROM pelanggan WHERE id_user = ? LIMIT 1"
	err := r.db.Where("id_user = ?", userID).First(&pelanggan).Error
	if err != nil {
		return nil, err // Kalau profil pelanggan tidak ketemu
	}

	return &pelanggan, nil // Kalau ketemu, kembalikan data
}