package model

import "time"

type Order struct {
	IDOrder             uint      `gorm:"primaryKey;autoIncrement;column:id_order"`
	PaketLayananID      uint      `gorm:"not null;column:id_paket_layanan"`
	PelangganID         uint      `gorm:"not null;column:id_pelanggan"` // Sesuai kesepakatan: id_pelanggan (bukan id_customer)
	AlamatPengambilanID uint      `gorm:"not null;column:id_alamat_pengambilan"`
	AlamatPenyerahanID  uint      `gorm:"not null;column:id_alamat_penyerahan"`
	ParfumID            uint      `gorm:"not null;column:id_parfum"`
	LayananID           uint      `gorm:"not null;column:id_layanan"`
	KaryawanID          *uint     `gorm:"column:id_karyawan"` // Pake pointer (*) karena kurir mungkin belum di-assign saat order masuk

	KeteranganLokasi string    `gorm:"type:text;column:keterangan_lokasi"`
	TglPesanan       time.Time `gorm:"type:timestamp;column:tgl_pesanan;default:CURRENT_TIMESTAMP"`
	JadwalPickup     time.Time `gorm:"type:timestamp;column:jadwal_pickup"`
	TipeLogistik     string    `gorm:"type:varchar(50);column:tipe_logistik"`
	HargaSaatIni     float64   `gorm:"type:numeric;column:harga_saat_ini"`
	Kuantitas        float64   `gorm:"type:numeric;column:kuantitas"`
	TotalBayar       float64   `gorm:"type:numeric;column:total_bayar"`
	CatatanOrder     string    `gorm:"type:text;column:catatan_order"`

	// --- Relasi GORM ---
	PaketLayanan      PaketLayanan `gorm:"foreignKey:PaketLayananID"`
	Pelanggan         Pelanggan    `gorm:"foreignKey:PelangganID"`
	AlamatPengambilan Alamat       `gorm:"foreignKey:AlamatPengambilanID"`
	AlamatPenyerahan  Alamat       `gorm:"foreignKey:AlamatPenyerahanID"`
	Parfum            Parfum       `gorm:"foreignKey:ParfumID"`
	Layanan           Layanan      `gorm:"foreignKey:LayananID"`
	Karyawan          Karyawan     `gorm:"foreignKey:KaryawanID"`
}

func (Order) TableName() string {
	return "order"
}