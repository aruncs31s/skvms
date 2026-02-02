package dto

type EssentialReadingRequest struct {
	Voltage float64 `json:"voltage" binding:"required"`
	Current float64 `json:"current" `
}
