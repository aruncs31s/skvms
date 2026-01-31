package dto

type GenericDropdown struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type VersionResponse struct {
	ID        uint             `json:"ID"`
	Version   string           `json:"Version"`
	CreatedAt string           `json:"CreatedAt"`
	UpdatedAt string           `json:"UpdatedAt"`
	Features  []FeatureResponse `json:"Features"`
}

type FeatureResponse struct {
	ID          uint   `json:"ID"`
	VersionID   uint   `json:"VersionID"`
	FeatureName string `json:"FeatureName"`
	Enabled     bool   `json:"Enabled"`
}
