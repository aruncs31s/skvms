package dto

import "time"

type EssentialReadingRequest struct {
	Voltage float64 `json:"voltage" binding:"required"`
	Current float64 `json:"current"`
}
type ProgressiveVoltageAndCurrent struct {
	Voltage    float64 `json:"voltage"`
	AvgVoltage float64 `json:"avg_voltage"`
	Current    float64 `json:"current"`
	AvgCurrent float64 `json:"avg_current"`
}
type ReadingsResponse struct {
	Voltage   float64   `json:"voltage"`
	Current   float64   `json:"current"`
	CreatedAt time.Time `json:"created_at"`
}
