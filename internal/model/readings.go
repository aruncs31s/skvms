package model

type Reading struct {
	ID        uint    `gorm:"column:id;primaryKey;autoIncrement"`
	DeviceID  uint    `gorm:"column:device_id;index;not null"`
	Voltage   float64 `gorm:"column:voltage"`
	Current   float64 `gorm:"column:current"`
	Timestamp int64   `gorm:"column:timestamp;index;not null"`
}

func (Reading) TableName() string {
	return "readings"
}
