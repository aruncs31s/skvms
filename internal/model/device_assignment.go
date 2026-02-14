package model

type DeviceAssignment struct {
	ID         uint `gorm:"column:id;primaryKey;autoIncrement"`
	LocationID uint `gorm:"column:location_id;index;not null"`
	DeviceID   uint `gorm:"column:device_id;index;not null"`
}

func (DeviceAssignment) TableName() string {
	return "device_assignment"
}
