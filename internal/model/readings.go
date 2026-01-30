package model

type Reading struct {
	ID        uint    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DeviceID  uint    `gorm:"column:device_id;index;not null" json:"device_id"`
	Voltage   float64 `gorm:"column:voltage" json:"voltage"`
	Current   float64 `gorm:"column:current" json:"current"`
	Timestamp int64   `gorm:"column:timestamp;index;not null" json:"timestamp"`
}

func (Reading) TableName() string {
	return "readings"
}
