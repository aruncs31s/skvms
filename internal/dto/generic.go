package dto

import "time"

type GenericDropdown struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type AdvancedGenericDropdown struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	DeviceType string    `json:"device_type"`
}
type GenericDropdownWithFeatures struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Features map[string]bool `json:"features"`
}

type VersionResponse struct {
	ID        uint              `json:"ID"`
	Version   string            `json:"Version"`
	CreatedAt string            `json:"CreatedAt"`
	UpdatedAt string            `json:"UpdatedAt"`
	Features  []FeatureResponse `json:"Features"`
}

type FeatureResponse struct {
	ID          uint   `json:"ID"`
	FeatureName string `json:"FeatureName"`
	Enabled     bool   `json:"Enabled"`
}

type AuditLogView struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	IPAddress string `json:"ip_address"`
	DeviceID  *uint  `json:"device_id,omitempty"`
	CreatedAt string `json:"created_at"`
}
