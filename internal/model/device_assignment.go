package model

import "time"

type DeviceAssignment struct {
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement"`
	LocationID   uint       `gorm:"column:location_id;index;not null"`
	DeviceID     uint       `gorm:"column:device_id;index;not null"`
	AssignedAt   time.Time  `gorm:"column:assigned_at;not null"`
	UnassignedAt *time.Time `gorm:"column:unassigned_at"`
}

func (DeviceAssignment) TableName() string {
	return "device_assignment"
}
