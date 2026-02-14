package model

import "time"

type DeviceAction uint8

type AuditLog struct {
	ID        uint         `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint         `gorm:"column:user_id;index" `
	Username  string       `gorm:"column:username"`
	Action    DeviceAction `gorm:"column:action;index"`
	Details   string       `gorm:"column:details"`
	IPAddress string       `gorm:"column:ip_address"`
	DeviceID  *uint        `gorm:"column:device_id;index"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	Device    *Device      `gorm:"foreignKey:DeviceID"`
	User      *User        `gorm:"foreignKey:UserID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
