package model

type Settings struct {
	IDSetting    uint   `gorm:"primaryKey;autoIncrement;column:id_setting"`
	SettingKey   string `gorm:"type:varchar(100);not null;unique;column:setting_key"`
	SettingValue string `gorm:"type:text;not null;column:setting_value"`
}

func (Settings) TableName() string {
	return "settings"
}