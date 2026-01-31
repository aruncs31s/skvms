package model

// Device Can be a Sensor, Actuator, Gateway, etc.
// I use device for esp32 and sensors
type Device struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
	// 1 -  , 2 - Sensor
	Type      uint `gorm:"column:device_type"`
	VersionID uint `gorm:"column:version_id"`

	// 1= Active, 0 = Inactive
	State      int             `gorm:"column:state"`
	Details    DeviceDetails   `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Addresses  []DeviceAddress `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeviceType DeviceTypes     `gorm:"foreignKey:Type;references:ID"`
	Version    Version         `gorm:"foreignKey:VersionID;references:ID"`

	CreatedBy uint `gorm:"column:created_by"`
	UpdatedBy uint `gorm:"column:updated_by"`

	CreatedAt int64 `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime"`
}

func (Device) TableName() string {
	return "devices"
}
