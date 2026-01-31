package model

// Our software has versioning and this will
// help to keep which device is which
type Version struct {
	ID      uint   `gorm:"primaryKey"`
	Version string `gorm:"unique;not null"`

	Features []Feature `gorm:"foreignKey:VersionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Version) TableName() string {
	return "versions"
}

// A version can have multiple features
// Also all features are backwards compatible
type Feature struct {
	ID          uint   `gorm:"primaryKey"`
	VersionID   uint   `gorm:"not null;index"`
	FeatureName string `gorm:"unique;not null"`
	Enabled     bool   `gorm:"not null;default:false"`
}

func (Feature) TableName() string {
	return "features"
}
