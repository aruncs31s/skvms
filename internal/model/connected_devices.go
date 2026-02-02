package model

// There might be connected devices associated with a main device
// Like LED , Fan etc
type ConnectedDevice struct {
	ParentID    uint   `gorm:"column:parent_id;index;not null"`
	ChildID     uint   `gorm:"column:device_id;index;not null"`
	ChildDevice Device `gorm:"foreignKey:ChildID;references:ID"` // Add this for preloading
}

func (ConnectedDevice) TableName() string {
	return "connected_devices"
}
