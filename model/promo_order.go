package model

type PromoOrder struct {
	IDPromoOrder uint `gorm:"primaryKey;autoIncrement;column:id_promo_order"`
	PromoID      uint `gorm:"not null;column:id_promo"`
	OrderID      uint `gorm:"not null;column:id_order"`

	Promo Promo `gorm:"foreignKey:PromoID"`
	Order Order `gorm:"foreignKey:OrderID"`
}

func (PromoOrder) TableName() string {
	return "promo_order"
}