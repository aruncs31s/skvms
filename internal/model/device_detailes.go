package model

import "time"

type DeviceDetails struct {
	ID              uint       `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID        uint       `gorm:"column:device_id;index;not null"`
	IPAddress       string     `gorm:"column:ip_address;index"`
	MACAddress      string     `gorm:"column:mac_address;uniqueIndex"`
	FirmwareVersion string     `gorm:"column:firmware_version"`
	LastSeenAt      *time.Time `gorm:"column:last_seen_at"`
}

func (DeviceDetails) TableName() string {
	return "device_details"
}
