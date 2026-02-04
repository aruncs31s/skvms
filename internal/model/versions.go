package model

import "time"

// Our software has versioning and this will
// help to keep which device is which version

type Version struct {
	ID                uint      `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	Name              string    `json:"Version" gorm:"column:name;unique;not null"`
	DeviceID          uint      `json:"DeviceID" gorm:"column:device_id;not null;index"`
	PreviousVersionID *uint     `json:"PreviousVersionID,omitempty" gorm:"column:previous_version_id;index;default:null"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
	Features          []Feature `json:"Features" gorm:"many2many:version_features;"`
	PreviousVersion   *Version  `json:"PreviousVersion,omitempty" gorm:"foreignKey:PreviousVersionID"`
}

func (Version) TableName() string {
	return "versions"
}

// A version can have multiple features
// Also all features are backwards compatible
type Feature struct {
	ID          uint      `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	FeatureName string    `json:"FeatureName" gorm:"column:feature_name;unique;not null"`
	Enabled     bool      `json:"Enabled" gorm:"column:enabled;not null;default:false"`
	Versions    []Version `json:"Versions" gorm:"many2many:version_features;"`
}

func (Feature) TableName() string {
	return "features"
}

// Junction table for many-to-many relationship between versions and features
type VersionFeature struct {
	VersionID uint `json:"VersionID" gorm:"column:version_id;primaryKey;not null"`
	FeatureID uint `json:"FeatureID" gorm:"column:feature_id;primaryKey;not null"`
}

func (VersionFeature) TableName() string {
	return "version_features"
}
