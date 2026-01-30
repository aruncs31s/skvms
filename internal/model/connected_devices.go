package model

// There might be connected devices associated with a main device
// Like LED , Fan etc
type ConnectedDevice struct {
	ID         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID   uint   `gorm:"column:device_id;index;not null"`
	DeviceName string `gorm:"column:device_name"`
}

func (ConnectedDevice) TableName() string {
	return "connected_devices"
}
