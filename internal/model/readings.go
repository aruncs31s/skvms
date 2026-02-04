package model

import "time"

// The HardwareTypeVoltageMeter measures voltage and current readings
type Reading struct {
	ID       uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DeviceID uint    `gorm:"column:device_id;index;not null" json:"device_id"`
	Voltage  float64 `gorm:"column:voltage" json:"voltage"`
	Current  float64 `gorm:"column:current" json:"current"`

	// Switch to CreatedAt to use time.Time for better handling
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Device    Device    `gorm:"foreignKey:DeviceID;references:ID"`
}

func (Reading) TableName() string {
	return "readings"
}
