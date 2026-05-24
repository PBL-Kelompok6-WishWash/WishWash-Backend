package repository

import (
	"time"

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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		var refStatus model.ReferensiStatusLayanan
		err := tx.Where("id_layanan = ? AND urutan_tahap = ?", order.LayananID, 1).First(&refStatus).Error
		if err != nil {
			err = tx.Where("id_layanan = ?", order.LayananID).Order("urutan_tahap asc").First(&refStatus).Error
		}

		if err == nil {
			history := model.RiwayatStatusDetail{
				ReferensiStatusID: refStatus.IDReferensiStatus,
				OrderID:           order.IDOrder,
				KaryawanID:        nil,
				WaktuUpdate:       time.Now(),
			}
			if err := tx.Create(&history).Error; err != nil {
				return err
			}
		}

		// Preload relationships back into the order struct after successful creation
		err = tx.Preload("PaketLayanan").
			Preload("Pelanggan").
			Preload("AlamatPengambilan").
			Preload("AlamatPenyerahan").
			Preload("Parfum").
			Preload("Layanan.ReferensiStatus").
			Preload("Karyawan").
			Preload("RiwayatStatusDetail.ReferensiStatus").
			Preload("Pembayaran").
			Preload("PromoOrder.Promo").
			First(order, order.IDOrder).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *orderRepository) FindAllByPelangganID(pelangganID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("PaketLayanan").
		Preload("Pelanggan").
		Preload("AlamatPengambilan").
		Preload("AlamatPenyerahan").
		Preload("Parfum").
		Preload("Layanan.ReferensiStatus").
		Preload("Karyawan").
		Preload("RiwayatStatusDetail.ReferensiStatus").
		Preload("Pembayaran").
		Preload("PromoOrder.Promo").
		Where("id_pelanggan = ?", pelangganID).
		Order("id_order desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(idOrder uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("PaketLayanan").
		Preload("Pelanggan").
		Preload("AlamatPengambilan").
		Preload("AlamatPenyerahan").
		Preload("Parfum").
		Preload("Layanan.ReferensiStatus").
		Preload("Karyawan").
		Preload("RiwayatStatusDetail.ReferensiStatus").
		Preload("Pembayaran").
		Preload("PromoOrder.Promo").
		First(&order, idOrder).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
