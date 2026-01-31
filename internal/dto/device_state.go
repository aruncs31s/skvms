package dto

import "time"

type DeviceStateFilterResponse struct {
	DeviceStates []GenericDropdown `json:"device_states"`
	TotalCount   int               `json:"total_count"`
}

type DeviceStateFilterRequest struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
	DeviceID uint      `json:"device_id"`
	States   []uint    `json:"states"`
}

type DeviceStateHistoryViewResponse struct {
	History      []DeviceStateHistoryView `json:"history"`
	TotalRecords int                      `json:"total_records"`
}

type DeviceStateHistoryView struct {
	StateName    string `json:"state_name"`
	ActionCaused string `json:"action_caused"`
	ChangedAt    string `json:"changed_at"`
	ChangedBy    string `json:"changed_by"`
}
