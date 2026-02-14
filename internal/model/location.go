package model

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID uint `gorm:"primaryKey;column:id"`
	// Sometype of code that used to identify the location ,
	// Set by the admin while creating.
	Code      string         `gorm:"column:code;type:varchar(50);unique;not null"`
	Name      string         `gorm:"column:name;type:varchar(255);not null"`
	Latitude  float64        `gorm:"column:latitude;type:decimal(10,8)"`
	Longitude float64        `gorm:"column:longitude;type:decimal(11,8)"`
	State     string         `gorm:"column:state;type:varchar(100)"`
	City      string         `gorm:"column:city;type:varchar(100)"`
	PinCode   string         `gorm:"column:pin_code;type:varchar(20)"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Location) TableName() string {
	return "locations"
}
