package model

import "time"

type HardwareType uint8

const (
	HardwareTypeUnknown HardwareType = iota
	HardwareTypeMicroController
	HardwareTypeSingleBoardComputer
	HardwareTypeSensor
)

var HardwareTypeMap = map[HardwareType]string{
	HardwareTypeUnknown:             "Unknown",
	HardwareTypeMicroController:     "MicroController",
	HardwareTypeSingleBoardComputer: "SingleBoardComputer",
	HardwareTypeSensor:              "Sensor",
}

func (h HardwareType) CanControl() bool {
	switch h {
	case HardwareTypeMicroController, HardwareTypeSingleBoardComputer:
		return false
	default:
		return true
	}
}
func (h HardwareType) String() string {
	if name, exists := HardwareTypeMap[h]; exists {
		return name
	}
	return "Unknown"
}

type DeviceTypes struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:type_name;uniqueIndex"`

	// 0: Unkown , 1: MicroController , 2: SingleBoardComputer, 3: Sensors
	HardwareType HardwareType `gorm:"column:hardware_type"`

	CreatedBy uint `gorm:"column:created_by"`
	UpdatedBy uint `gorm:"column:updated_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DeviceTypes) TableName() string {
	return "device_types"
}
