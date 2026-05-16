package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type AlamatRepository interface {
	FindByPelangganID(pelangganID uint) (*model.Alamat, error)
}

type alamatRepository struct {
	db *gorm.DB
}

func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db}
}

func (r *alamatRepository) FindByPelangganID(pelangganID uint) (*model.Alamat, error) {
	var alamat model.Alamat
	err := r.db.Where("id_pelanggan = ?", pelangganID).First(&alamat).Error
	if err != nil {
		return nil, err
	}
	return &alamat, nil
}
