package repository

import (
	"log"

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
	UpdatePaketLayanan(layananID uint, pakets []model.PaketLayanan) error
	CheckIsUsed(id uint) (bool, error)
}

type layananRepository struct {
	db *gorm.DB
}

func NewLayananRepository(db *gorm.DB) LayananRepository {
	return &layananRepository{db}
}

func (r *layananRepository) FindAll() ([]model.Layanan, error) {
	var layanans []model.Layanan
	// Preload referensi status dan paket layanan
	err := r.db.Preload("ReferensiStatus", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tahap ASC")
	}).Preload("PaketLayanan").Order("id_layanan ASC").Find(&layanans).Error
	return layanans, err
}

func (r *layananRepository) FindByID(id uint) (*model.Layanan, error) {
	var layanan model.Layanan
	err := r.db.Preload("ReferensiStatus", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tahap ASC")
	}).Preload("PaketLayanan").First(&layanan, id).Error
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
	// Updates the parent fields, ignoring associations to avoid FK constraint errors during save
	return r.db.Omit("ReferensiStatus", "PaketLayanan").Save(layanan).Error
}

func (r *layananRepository) UpdateStatusLayanan(layananID uint, statuses []model.ReferensiStatusLayanan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Dapatkan status yang ada saat ini di DB
		var existing []model.ReferensiStatusLayanan
		if err := tx.Where("id_layanan = ?", layananID).Find(&existing).Error; err != nil {
			return err
		}

		existingMap := make(map[string]model.ReferensiStatusLayanan)
		for _, s := range existing {
			existingMap[s.NamaStatus] = s
		}

		keptNames := make(map[string]bool)

		// 2. Loop data input untuk create atau update
		for _, s := range statuses {
			keptNames[s.NamaStatus] = true
			if ext, found := existingMap[s.NamaStatus]; found {
				// Update urutan jika ada perubahan
				ext.UrutanTahap = s.UrutanTahap
				if err := tx.Save(&ext).Error; err != nil {
					return err
				}
			} else {
				// Tambah baru
				s.LayananID = layananID
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
			}
		}

		// 3. Hapus status yang tidak ada lagi di input
		for _, s := range existing {
			if !keptNames[s.NamaStatus] {
				// Jika gagal karena constraint, biarkan saja agar data riwayat tidak terputus
				if err := tx.Delete(&s).Error; err != nil {
					log.Printf("⚠️ WARNING: Gagal menghapus referensi_status_layanan ID %d (%s) karena masih digunakan oleh riwayat order: %v\n", s.IDReferensiStatus, s.NamaStatus, err)
				}
			}
		}
		return nil
	})
}

func (r *layananRepository) UpdatePaketLayanan(layananID uint, pakets []model.PaketLayanan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Dapatkan paket yang ada saat ini di DB
		var existing []model.PaketLayanan
		if err := tx.Where("id_layanan = ?", layananID).Find(&existing).Error; err != nil {
			return err
		}

		existingMap := make(map[string]model.PaketLayanan)
		for _, p := range existing {
			existingMap[p.NamaPaket] = p
		}

		keptNames := make(map[string]bool)

		// 2. Loop data input untuk create atau update
		for _, p := range pakets {
			keptNames[p.NamaPaket] = true
			if ext, found := existingMap[p.NamaPaket]; found {
				// Update detail paket
				ext.DurasiJam = p.DurasiJam
				ext.BiayaTambahan = p.BiayaTambahan
				if err := tx.Save(&ext).Error; err != nil {
					return err
				}
			} else {
				// Tambah baru
				p.LayananID = layananID
				if err := tx.Create(&p).Error; err != nil {
					return err
				}
			}
		}

		// 3. Hapus paket yang tidak ada lagi di input
		for _, p := range existing {
			if !keptNames[p.NamaPaket] {
				// Jika gagal karena constraint, biarkan saja agar data order lama tidak terputus
				if err := tx.Delete(&p).Error; err != nil {
					log.Printf("⚠️ WARNING: Gagal menghapus paket_layanan ID %d (%s) karena masih digunakan oleh order: %v\n", p.IDPaketLayanan, p.NamaPaket, err)
				}
			}
		}
		return nil
	})
}

func (r *layananRepository) Delete(id uint) error {
	// Explicitly delete child tables first to avoid FK constraint error
	r.db.Where("id_layanan = ?", id).Delete(&model.ReferensiStatusLayanan{})
	r.db.Where("id_layanan = ?", id).Delete(&model.PaketLayanan{})
	return r.db.Delete(&model.Layanan{}, id).Error
}

func (r *layananRepository) CheckIsUsed(id uint) (bool, error) {
	var count int64
	err := r.db.Table("order").Where("id_layanan = ?", id).Count(&count).Error
	return count > 0, err
}
