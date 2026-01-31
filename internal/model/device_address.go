package model

type DeviceAddress struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID uint   `gorm:"column:device_id;index;not null"`
	Address  string `gorm:"column:address"`
	City     string `gorm:"column:city"`
}

func (DeviceAddress) TableName() string {
	return "device_address"
}
