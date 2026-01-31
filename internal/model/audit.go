package model

import "time"

type DeviceAction uint8

type AuditLog struct {
	ID        uint         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint         `gorm:"column:user_id;index" json:"user_id"`
	Username  string       `gorm:"column:username" json:"username"`
	Action    DeviceAction `gorm:"column:action;index" json:"action"`
	Details   string       `gorm:"column:details" json:"details"`
	IPAddress string       `gorm:"column:ip_address" json:"ip_address"`
	DeviceID  *uint        `gorm:"column:device_id;index" json:"device_id,omitempty"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Device    *Device      `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
