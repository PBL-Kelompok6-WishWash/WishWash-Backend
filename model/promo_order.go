package model

type PromoOrder struct {
	IDPromoOrder uint `gorm:"primaryKey;autoIncrement;column:id_promo_order" json:"id_promo_order"`
	PromoID      uint `gorm:"not null;column:id_promo" json:"id_promo"`
	OrderID      uint `gorm:"not null;column:id_order" json:"id_order"`

	Promo Promo `gorm:"foreignKey:PromoID" json:"Promo"`
	Order Order `gorm:"foreignKey:OrderID" json:"Order"`
}

func (PromoOrder) TableName() string {
	return "promo_order"
}