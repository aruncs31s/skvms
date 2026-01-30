package model

// Our software has versioning and this will
// help to keep which device is which
type Versions struct {
	ID      uint   `gorm:"primaryKey"`
	Version string `gorm:"unique;not null"`
}

func (Versions) TableName() string {
	return "versions"
}

// A version can have multiple features
// Also all features are backwards compatible
type Features struct {
	ID          uint   `gorm:"primaryKey"`
	VersionID   uint   `gorm:"not null;index"`
	FeatureName string `gorm:"unique;not null"`
	Enabled     bool   `gorm:"not null;default:false"`
}

func (Features) TableName() string {
	return "features"
}
