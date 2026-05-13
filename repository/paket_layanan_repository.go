package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type PaketLayananRepository interface {
	FindAll() ([]model.PaketLayanan, error)
	FindByID(id uint) (*model.PaketLayanan, error)
	Create(paket *model.PaketLayanan) error
	Update(paket *model.PaketLayanan) error
	Delete(id uint) error
}

type paketLayananRepository struct {
	db *gorm.DB
}

func NewPaketLayananRepository(db *gorm.DB) PaketLayananRepository {
	return &paketLayananRepository{db}
}

func (r *paketLayananRepository) FindAll() ([]model.PaketLayanan, error) {
	var pakets []model.PaketLayanan
	err := r.db.Find(&pakets).Error
	return pakets, err
}

func (r *paketLayananRepository) FindByID(id uint) (*model.PaketLayanan, error) {
	var paket model.PaketLayanan
	err := r.db.First(&paket, id).Error
	if err != nil {
		return nil, err
	}
	return &paket, nil
}

func (r *paketLayananRepository) Create(paket *model.PaketLayanan) error {
	return r.db.Create(paket).Error
}

func (r *paketLayananRepository) Update(paket *model.PaketLayanan) error {
	return r.db.Save(paket).Error
}

func (r *paketLayananRepository) Delete(id uint) error {
	return r.db.Delete(&model.PaketLayanan{}, id).Error
}
