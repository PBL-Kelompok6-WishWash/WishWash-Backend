package repository

import (
	"github.com/PBL-Kelompok6-WishWash/backend/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindAllByPelangganID(pelangganID uint) ([]model.Order, error)
	FindByID(idOrder uint) (*model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindAllByPelangganID(pelangganID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("PaketLayanan").
		Preload("AlamatPengambilan").
		Preload("AlamatPenyerahan").
		Preload("Parfum").
		Preload("Layanan").
		Preload("Karyawan").
		Where("id_pelanggan = ?", pelangganID).
		Order("id_order desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(idOrder uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("PaketLayanan").
		Preload("AlamatPengambilan").
		Preload("AlamatPenyerahan").
		Preload("Parfum").
		Preload("Layanan").
		Preload("Karyawan").
		First(&order, idOrder).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
