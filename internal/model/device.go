package model

type Device struct {
	ID        uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string `gorm:"column:name"`
	Type      uint   `gorm:"column:device_type"`
	VersionID uint   `gorm:"column:version_id"`
	CreatedBy uint   `gorm:"column:created_by"`
	UpdatedBy uint   `gorm:"column:updated_by"`
	// 1= Active, 0 = Inactive
	State      int             `gorm:"column:state"`
	Details    DeviceDetails   `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Addresses  []DeviceAddress `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DeviceType DeviceTypes     `gorm:"foreignKey:Type;references:ID"`
	Version    Version         `gorm:"foreignKey:VersionID;references:ID"`
}

func (Device) TableName() string {
	return "devices"
}
