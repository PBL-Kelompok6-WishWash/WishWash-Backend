package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type AlamatRepository interface {
	FindByPelangganID(pelangganID uint) (*model.Alamat, error)
	FindAllByPelangganID(pelangganID uint) ([]model.Alamat, error)
	Create(alamat *model.Alamat) error
	Delete(idAlamat uint, pelangganID uint) error
	SetPrimary(idAlamat uint, pelangganID uint) error
	Update(alamat *model.Alamat) error
}

type alamatRepository struct {
	db *gorm.DB
}

func NewAlamatRepository(db *gorm.DB) AlamatRepository {
	return &alamatRepository{db}
}

func (r *alamatRepository) FindByPelangganID(pelangganID uint) (*model.Alamat, error) {
	var alamat model.Alamat
	err := r.db.Where("id_pelanggan = ?", pelangganID).Order("is_primary desc, id_alamat desc").First(&alamat).Error
	if err != nil {
		return nil, err
	}
	return &alamat, nil
}

func (r *alamatRepository) FindAllByPelangganID(pelangganID uint) ([]model.Alamat, error) {
	var alamats []model.Alamat
	err := r.db.Where("id_pelanggan = ?", pelangganID).Order("is_primary desc, id_alamat desc").Find(&alamats).Error
	return alamats, err
}

func (r *alamatRepository) Create(alamat *model.Alamat) error {
	// Jika alamat ini adalah primary, set semua alamat lain jadi tidak primary
	if alamat.IsPrimary {
		r.db.Model(&model.Alamat{}).Where("id_pelanggan = ?", alamat.PelangganID).Update("is_primary", false)
	} else {
		// Jika belum ada alamat sama sekali, jadikan primary
		var count int64
		r.db.Model(&model.Alamat{}).Where("id_pelanggan = ?", alamat.PelangganID).Count(&count)
		if count == 0 {
			alamat.IsPrimary = true
		}
	}
	return r.db.Create(alamat).Error
}

func (r *alamatRepository) Delete(idAlamat uint, pelangganID uint) error {
	return r.db.Where("id_alamat = ? AND id_pelanggan = ?", idAlamat, pelangganID).Delete(&model.Alamat{}).Error
}

func (r *alamatRepository) SetPrimary(idAlamat uint, pelangganID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Set semua false
		if err := tx.Model(&model.Alamat{}).Where("id_pelanggan = ?", pelangganID).Update("is_primary", false).Error; err != nil {
			return err
		}
		// Set target true
		if err := tx.Model(&model.Alamat{}).Where("id_alamat = ? AND id_pelanggan = ?", idAlamat, pelangganID).Update("is_primary", true).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *alamatRepository) Update(alamat *model.Alamat) error {
	return r.db.Save(alamat).Error
}
