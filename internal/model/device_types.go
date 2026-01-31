package model

import "time"

type DeviceTypes struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:type_name;uniqueIndex"`

	// 0: Unkown , 1: MicroController , 2: SingleBoardComputer, 3: Sensors

	HardwareType uint8 `gorm:"column:hardware_type"`

	CreatedBy uint `gorm:"column:created_by"`
	UpdatedBy uint `gorm:"column:updated_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DeviceTypes) TableName() string {
	return "device_types"
}
