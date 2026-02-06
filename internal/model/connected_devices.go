package model

import "gorm.io/gorm"

// There might be connected devices associated with a main device
// Like LED , Fan etc
type ConnectedDevice struct {
	ParentID    uint           `gorm:"column:parent_id;index;not null"`
	ChildID     uint           `gorm:"column:child_id;index;not null"`
	Deleted     gorm.DeletedAt `gorm:"index"`
	ChildDevice Device         `gorm:"foreignKey:ChildID;references:ID"` // Add this for preloading
}

func (ConnectedDevice) TableName() string {
	return "connected_devices"
}

type ConnectedDeviceReadings struct {
	ParentDevice         uint    `gorm:"column:parent_id"`
	ChildDevice          uint    `gorm:"column:child_id"`
	Voltage              float64 `gorm:"column:voltage"`
	ChargingCurrent      float64 `gorm:"column:current"`
	AvgVoltage           float64 `gorm:"column:avg_voltage"`
	AvgChargingCurrent   float64 `gorm:"column:avg_current"`
	BatteryRemainingTime float64 `gorm:"column:estimated_remaining_hours"`
	CreatedAt            string  `gorm:"column:created_at"`
}
