package dto

type DeviceView struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	IPAddress       string `json:"ip_address"`
	MACAddress      string `json:"mac_address"`
	FirmwareVersion string `json:"firmware_version"`
	VersionID       uint   `json:"version_id"`
	Address         string `json:"address"`
	City            string `json:"city"`
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

type UpdateDeviceRequest struct {
	Name              string `json:"name" binding:"required"`
	Type              uint   `json:"type" binding:"required"`
	IPAddress         string `json:"ip_address" binding:"required"`
	MACAddress        string `json:"mac_address" binding:"required"`
	FirmwareVersionID uint   `json:"firmware_version_id" binding:"required"`
	Address           string `json:"address"`
	City              string `json:"city"`
}
