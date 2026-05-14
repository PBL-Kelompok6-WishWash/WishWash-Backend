package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type PromoRepository interface {
	FindAll() ([]model.Promo, error)
	FindByID(id int) (model.Promo, error)
	Create(promo model.Promo) (model.Promo, error)
	Update(promo model.Promo) (model.Promo, error)
	Delete(promo model.Promo) error
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) *promoRepository {
	return &promoRepository{db}
}

func (r *promoRepository) FindAll() ([]model.Promo, error) {
	var promos []model.Promo
	err := r.db.Order("id_promo asc").Find(&promos).Error
	return promos, err
}

func (r *promoRepository) FindByID(id int) (model.Promo, error) {
	var promo model.Promo
	err := r.db.Where("id_promo = ?", id).First(&promo).Error
	return promo, err
}

func (r *promoRepository) Create(promo model.Promo) (model.Promo, error) {
	err := r.db.Create(&promo).Error
	return promo, err
}

func (r *promoRepository) Update(promo model.Promo) (model.Promo, error) {
	err := r.db.Save(&promo).Error
	return promo, err
}

func (r *promoRepository) Delete(promo model.Promo) error {
	return r.db.Delete(&promo).Error
}
