package model

import "time"

type DeviceView struct {
	ID              uint         `gorm:"column:id"`
	Name            string       `gorm:"column:name"`
	Type            string       `gorm:"column:type"`
	HardwareType    HardwareType `gorm:"column:hardware_type"`
	IPAddress       string       `gorm:"column:ip_address"`
	MACAddress      string       `gorm:"column:mac_address"`
	FirmwareVersion string       `gorm:"column:firmware_version"`
	Address         string       `gorm:"column:address"`
	City            string       `gorm:"column:city"`
	DeviceState     string       `gorm:"column:current_state"`
}
type MicrocontrollerDeviceView struct {
	ID              uint   `gorm:"column:id"`
	ParentID        uint   `gorm:"column:parent_id"`
	Name            string `gorm:"column:name"`
	Type            string `gorm:"column:type"`
	IPAddress       string `gorm:"column:ip_address"`
	MACAddress      string `gorm:"column:mac_address"`
	FirmwareVersion string `gorm:"column:firmware_version"`
	DeviceState     string `gorm:"column:current_state"` // Status
	UsedBy          string `gorm:"column:used_by"`
}
type MicrocontrollerStatsView struct {
	TotalDevices   int64 `json:"total_devices"`
	OnlineDevices  int64 `json:"online_devices"`
	OfflineDevices int64 `json:"offline_devices"`
}
type DeviceDetails struct {
	ID         uint       `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID   uint       `gorm:"column:device_id;index;not null"`
	IPAddress  string     `gorm:"column:ip_address;index"`
	MACAddress string     `gorm:"column:mac_address"`
	LastSeenAt *time.Time `gorm:"column:last_seen_at"`
}

func (DeviceDetails) TableName() string {
	return "device_details"
}
