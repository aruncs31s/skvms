package model

import "time"

// Device Can be a Sensor, Actuator, Gateway, etc.
// I use device for esp32 and sensors
type Device struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`

	// 1 -  , 2 - Sensor
	DeviceTypeID uint `gorm:"column:device_type"`
	VersionID    uint `gorm:"column:version_id"`

	// 1= Active, 0 = Inactive , 2 = Maintenance, 3 = Decommissioned
	// Also FK to DeviceState.ID
	CurrentState int `gorm:"column:device_state"`

	Details    DeviceDetails   `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Addresses  []DeviceAddress `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeviceType DeviceTypes     `gorm:"foreignKey:Type;references:ID"`
	Version    Version         `gorm:"foreignKey:VersionID;references:ID"`

	CreatedBy uint `gorm:"column:created_by"`
	UpdatedBy uint `gorm:"column:updated_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Device) TableName() string {
	return "devices"
}

// Possible States for a device
// Different types of devices can have different states
type DeviceState struct {
	ID           int       `gorm:"column:id;primaryKey"`
	Name         string    `gorm:"column:name"`
	DeviceTypeID uint      `gorm:"column:type"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DeviceState) TableName() string {
	return "device_states"
}

type DeviceStateHistory struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID  uint      `gorm:"column:device_id"`
	StateID   int       `gorm:"column:state_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DeviceStateHistory) TableName() string {
	return "device_state_history"
}
