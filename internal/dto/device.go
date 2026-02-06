package dto

import "github.com/aruncs31s/skvms/internal/model"

type DeviceView struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Like For Current Sensor We Have CT and Hall Effect
	Type            string             `json:"type"`
	HardwareType    model.HardwareType `json:"hardware_type"`
	Status          string             `json:"status"`
	IPAddress       string             `json:"ip_address"`
	MACAddress      string             `json:"mac_address"`
	FirmwareVersion string             `json:"firmware_version"`
	Address         string             `json:"address"`
	City            string             `json:"city"`
}
type MicrocontrollerDeviceView struct {
	ID               uint               `json:"id"`
	ParentID         *uint              `json:"parent_id,omitempty"`
	Name             string             `json:"name"`
	IP               string             `json:"ip_address"`
	MAC              string             `json:"mac_address"`
	Status           string             `json:"status"`
	UsedBy           *string            `json:"used_by"` // Or Connected To this perticullar Solar Device
	FirmwareVersion  string             `json:"firmware_version"`
	ConncetedSensors []SensorDeviceView `json:"connected_sensors"`
}

type MicrocontrollerStatsView struct {
	TotalMicrocontrollers   int64  `json:"total_microcontrollers"`
	OnlineMicrocontrollers  int64  `json:"online_microcontrollers"`
	OfflineMicrocontrollers int64  `json:"offline_microcontrollers"`
	LatestVersion           string `json:"latest_version"`
}

type SensorDeviceView struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type CreateDeviceRequest struct {
	Name              string `json:"name" binding:"required"`
	Type              uint   `json:"type" binding:"required"`
	IPAddress         string `json:"ip_address" `
	MACAddress        string `json:"mac_address"`
	FirmwareVersionID uint   `json:"firmware_version_id"`
	Address           string `json:"address"`
	City              string `json:"city"`
}

type CreateConnectedDeviceRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       uint   `json:"type" binding:"required"`
	IPAddress  string `json:"ip_address" `
	MACAddress string `json:"mac_address"`
}

type UpdateDeviceRequest struct {
	Name              *string `json:"name,omitempty"`
	Type              *uint   `json:"type,omitempty"`
	IPAddress         *string `json:"ip_address,omitempty"`
	MACAddress        *string `json:"mac_address,omitempty"`
	FirmwareVersionID *uint   `json:"firmware_version_id,omitempty"`
	Address           *string `json:"address,omitempty"`
	City              *string `json:"city,omitempty"`
}

type FullUpdateDeviceRequest struct {
	Name              string `json:"name" binding:"required"`
	Type              uint   `json:"type" binding:"required"`
	IPAddress         string `json:"ip_address" binding:"required"`
	MACAddress        string `json:"mac_address" binding:"required"`
	FirmwareVersionID uint   `json:"firmware_version_id" binding:"required"`
	Address           string `json:"address" binding:"required"`
	City              string `json:"city" binding:"required"`
	CurrentState      uint   `json:"current_state" binding:"required"`
}

type DeviceStateView struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type CreateDeviceStateRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateDeviceStateRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddConnectedDeviceRequest struct {
	ChildID uint `json:"child_id" binding:"required"`
}
