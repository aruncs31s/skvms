package model

import "time"

type Device struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `gorm:"column:name"`
	Type      string `gorm:"column:device_type"`
	CreatedBy uint   `gorm:"column:created_by"`
	UpdatedBy uint   `gorm:"column:updated_by"`

	// 1= Active, 0 = Inactive
	State     int             `gorm:"column:state"`
	Details   DeviceDetails   `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Addresses []DeviceAddress `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DeviceDetails struct {
	ID              uint       `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID        uint       `gorm:"column:device_id;index;not null"`
	IPAddress       string     `gorm:"column:ip_address;index"`
	MACAddress      string     `gorm:"column:mac_address;uniqueIndex"`
	FirmwareVersion string     `gorm:"column:firmware_version"`
	LastSeenAt      *time.Time `gorm:"column:last_seen_at"`
}

type DeviceAddress struct {
	ID       uint   `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID uint   `gorm:"column:device_id;index;not null"`
	Address  string `gorm:"column:address"`
	City     string `gorm:"column:city"`
}

func (Device) TableName() string {
	return "devices"
}

func (DeviceDetails) TableName() string {
	return "device_details"
}

func (DeviceAddress) TableName() string {
	return "device_address"
}
