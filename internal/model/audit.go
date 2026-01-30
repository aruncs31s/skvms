package model

import "time"

type AuditLog struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"column:user_id;index" json:"user_id"`
	Username  string    `gorm:"column:username" json:"username"`
	Action    string    `gorm:"column:action;index" json:"action"`
	Details   string    `gorm:"column:details" json:"details"`
	IPAddress string    `gorm:"column:ip_address" json:"ip_address"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
