package model

type DeviceTypes struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:type_name;uniqueIndex"`
}

func (DeviceTypes) TableName() string {
	return "device_types"
}
