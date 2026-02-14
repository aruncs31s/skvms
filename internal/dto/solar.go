package dto

type SolarDeviceView struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	ChargingCurrent   float64 `json:"charging_current"`
	BatteryVoltage    float64 `json:"battery_voltage"`
	RemainingTime     float64 `json:"remaining_time"`
	LedStatus         string  `json:"led_status"`
	ConnectedDeviceIP string  `json:"connected_device_ip"`
	Address           string  `json:"address"`
	City              string  `json:"city"`
	Status            string  `json:"status"`
}
type CreateSolarDeviceDTO struct {
	Name                       string `json:"name" binding:"required"`
	DeviceTypeID               uint   `json:"device_type_id" binding:"required"`
	ConnectedMicroControllerID *uint  `json:"connected_microcontroller_id,omitempty"`
	Address                    string `json:"address" binding:"required"`
	City                       string `json:"city" binding:"required"`
}
