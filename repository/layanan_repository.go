package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type LayananRepository interface {
	FindAll() ([]model.Layanan, error)
	FindByID(id uint) (*model.Layanan, error)
	Create(layanan *model.Layanan) error
	Update(layanan *model.Layanan) error
	Delete(id uint) error
	UpdateStatusLayanan(layananID uint, statuses []model.ReferensiStatusLayanan) error
}

type layananRepository struct {
	db *gorm.DB
}

func NewLayananRepository(db *gorm.DB) LayananRepository {
	return &layananRepository{db}
}

func (r *layananRepository) FindAll() ([]model.Layanan, error) {
	var layanans []model.Layanan
	// Preload referensi status, order by urutan_tahap
	err := r.db.Preload("ReferensiStatus", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tahap ASC")
	}).Find(&layanans).Error
	return layanans, err
}

func (r *layananRepository) FindByID(id uint) (*model.Layanan, error) {
	var layanan model.Layanan
	err := r.db.Preload("ReferensiStatus", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tahap ASC")
	}).First(&layanan, id).Error
	if err != nil {
		return nil, err
	}
	return &layanan, nil
}

func (r *layananRepository) Create(layanan *model.Layanan) error {
	// GORM will automatically insert the nested ReferensiStatus because of the relationship
	return r.db.Create(layanan).Error
}

func (r *layananRepository) Update(layanan *model.Layanan) error {
	// Updates the parent fields
	return r.db.Save(layanan).Error
}

func (r *layananRepository) UpdateStatusLayanan(layananID uint, statuses []model.ReferensiStatusLayanan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete all existing statuses for this Layanan
		if err := tx.Where("id_layanan = ?", layananID).Delete(&model.ReferensiStatusLayanan{}).Error; err != nil {
			return err
		}

		// 2. Insert the new ones if there are any
		if len(statuses) > 0 {
			if err := tx.Create(&statuses).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *layananRepository) Delete(id uint) error {
	// Explicitly delete child statuses first to avoid FK constraint error if DB schema wasn't dropped
	r.db.Where("id_layanan = ?", id).Delete(&model.ReferensiStatusLayanan{})
	return r.db.Delete(&model.Layanan{}, id).Error
}
