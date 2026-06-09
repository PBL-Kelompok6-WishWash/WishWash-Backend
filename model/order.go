package model

import "time"

type Order struct {
	IDOrder             uint      `gorm:"primaryKey;autoIncrement;column:id_order" json:"id_order"`
	KodeOrder           string    `gorm:"type:varchar(50);column:kode_order;unique" json:"kode_order"`
	PaketLayananID      *uint     `gorm:"column:id_paket_layanan" json:"id_paket_layanan"` // Bisa null jika layanan tidak punya paket
	PelangganID         uint      `gorm:"not null;column:id_pelanggan" json:"id_pelanggan"` // Sesuai kesepakatan: id_pelanggan (bukan id_customer)
	AlamatPengambilanID *uint     `gorm:"column:id_alamat_pengambilan" json:"id_alamat_pengambilan"`
	AlamatPenyerahanID  *uint     `gorm:"column:id_alamat_penyerahan" json:"id_alamat_penyerahan"`
	ParfumID            uint      `gorm:"not null;column:id_parfum" json:"id_parfum"`
	LayananID           uint      `gorm:"not null;column:id_layanan" json:"id_layanan"`
	KaryawanID          *uint     `gorm:"column:id_karyawan" json:"id_karyawan"` // Pake pointer (*) karena kurir mungkin belum di-assign saat order masuk

	KeteranganLokasi string     `gorm:"type:text;column:keterangan_lokasi" json:"keterangan_lokasi"`
	TglPesanan       time.Time  `gorm:"type:timestamp;column:tgl_pesanan;default:CURRENT_TIMESTAMP" json:"tgl_pesanan"`
	JadwalPickup     *time.Time `gorm:"type:timestamp;column:jadwal_pickup" json:"jadwal_pickup"`
	TipeLogistik     string    `gorm:"type:varchar(50);column:tipe_logistik" json:"tipe_logistik"`
	HargaSaatIni     float64   `gorm:"type:numeric;column:harga_saat_ini" json:"harga_saat_ini"`
	Kuantitas        float64   `gorm:"type:numeric;column:kuantitas" json:"kuantitas"`
	TotalBayar       float64   `gorm:"type:numeric;column:total_bayar" json:"total_bayar"`
	CatatanOrder     string    `gorm:"type:text;column:catatan_order" json:"catatan_order"`
	IsCourierOnWay   bool      `gorm:"column:is_courier_on_way;default:false" json:"is_courier_on_way"`
	CourierLatitude  string    `gorm:"type:varchar(100);column:courier_latitude" json:"courier_latitude"`
	CourierLongitude string    `gorm:"type:varchar(100);column:courier_longitude" json:"courier_longitude"`

	// --- Relasi GORM ---
	PaketLayanan        *PaketLayanan         `gorm:"foreignKey:PaketLayananID" json:"PaketLayanan"`
	Pelanggan           Pelanggan             `gorm:"foreignKey:PelangganID" json:"Pelanggan"`
	AlamatPengambilan   *Alamat                `gorm:"foreignKey:AlamatPengambilanID" json:"AlamatPengambilan"`
	AlamatPenyerahan    *Alamat                `gorm:"foreignKey:AlamatPenyerahanID" json:"AlamatPenyerahan"`
	Parfum              Parfum                `gorm:"foreignKey:ParfumID" json:"Parfum"`
	Layanan             Layanan               `gorm:"foreignKey:LayananID" json:"Layanan"`
	Karyawan            *Karyawan              `gorm:"foreignKey:KaryawanID" json:"Karyawan"`
	RiwayatStatusDetail []RiwayatStatusDetail `gorm:"foreignKey:OrderID" json:"RiwayatStatusDetail"`
	Pembayaran          *Pembayaran           `gorm:"foreignKey:OrderID" json:"Pembayaran"`
	PromoOrder          []PromoOrder          `gorm:"foreignKey:OrderID" json:"PromoOrder"`
	Penilaian           *Penilaian            `gorm:"foreignKey:OrderID" json:"Penilaian"`
}

func (Order) TableName() string {
	return "order"
}