package dto

type DeviceView struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	IPAddress       string `json:"ip_address"`
	MACAddress      string `json:"mac_address"`
	FirmwareVersion string `json:"firmware_version"`
	VersionID       uint   `json:"version_id"`
	Address         string `json:"address"`
	City            string `json:"city"`
	DeviceState     string `json:"device_state"`
}

type CreateDeviceRequest struct {
	Name              string `json:"name" binding:"required"`
	UID               string `json:"uid" binding:"required"`
	Type              uint   `json:"type"`
	IPAddress         string `json:"ip_address" `
	MACAddress        string `json:"mac_address"`
	FirmwareVersionID uint   `json:"firmware_version_id"`
	Address           string `json:"address"`
	City              string `json:"city"`
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
