package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type ParfumRepository interface {
	FindAll() ([]model.Parfum, error)
	FindByID(id int) (model.Parfum, error)
	Create(parfum model.Parfum) (model.Parfum, error)
	Update(parfum model.Parfum) (model.Parfum, error)
	Delete(parfum model.Parfum) error
	CheckIsUsed(id int) (bool, error)
}

type parfumRepository struct {
	db *gorm.DB
}

func NewParfumRepository(db *gorm.DB) *parfumRepository {
	return &parfumRepository{db}
}

func (r *parfumRepository) FindAll() ([]model.Parfum, error) {
	var parfums []model.Parfum
	err := r.db.Find(&parfums).Error
	return parfums, err
}

func (r *parfumRepository) FindByID(id int) (model.Parfum, error) {
	var parfum model.Parfum
	err := r.db.Where("id_parfum = ?", id).First(&parfum).Error
	return parfum, err
}

func (r *parfumRepository) Create(parfum model.Parfum) (model.Parfum, error) {
	err := r.db.Create(&parfum).Error
	return parfum, err
}

func (r *parfumRepository) Update(parfum model.Parfum) (model.Parfum, error) {
	err := r.db.Save(&parfum).Error
	return parfum, err
}

func (r *parfumRepository) Delete(parfum model.Parfum) error {
	err := r.db.Delete(&parfum).Error
	return err
}

func (r *parfumRepository) CheckIsUsed(id int) (bool, error) {
	var count int64
	err := r.db.Table("order").Where("id_parfum = ?", id).Count(&count).Error
	return count > 0, err
}
