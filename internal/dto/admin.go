package dto

type AdminAggregateResponse struct {
	TotalUsers      int `json:"total_users"`
	TotalDevices    int `json:"total_devices"`
	ActiveDevices   int `json:"active_devices"`
	InactiveDevices int `json:"inactive_devices"`
	TotalReadings   int `json:"total_readings"`
	TotalAudits     int `json:"total_audits"`
}
