package model

import "time"

// Our software has versioning and this will
// help to keep which device is which version

type Version struct {
	ID                uint      `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	Version           string    `json:"Version" gorm:"column:version;unique;not null"`
	PreviousVersionID *uint     `json:"PreviousVersionID,omitempty" gorm:"column:previous_version_id;index;default:null"`
	CreatedAt         time.Time `json:"CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt"`
	Features          []Feature `json:"Features" gorm:"foreignKey:VersionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PreviousVersion   *Version  `json:"PreviousVersion,omitempty" gorm:"foreignKey:PreviousVersionID"`
}

func (Version) TableName() string {
	return "versions"
}

// A version can have multiple features
// Also all features are backwards compatible
type Feature struct {
	ID          uint   `json:"ID" gorm:"column:id;primaryKey;autoIncrement"`
	VersionID   uint   `json:"VersionID" gorm:"column:version_id;not null;index"`
	FeatureName string `json:"FeatureName" gorm:"column:feature_name;unique;not null"`
	Enabled     bool   `json:"Enabled" gorm:"column:enabled;not null;default:false"`
}

func (Feature) TableName() string {
	return "features"
}
