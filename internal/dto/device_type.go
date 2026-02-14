package dto

type CreateDeviceTypeRequest struct {
	Name         string `json:"name" binding:"required"`
	HardwareType uint   `json:"hardware_type" binding:"required"`
}
