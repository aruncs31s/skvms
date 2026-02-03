package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type VersionRepository interface {
	CreateVersion(ctx context.Context, version *model.Version) error
	GetAllVersions(ctx context.Context) ([]model.Version, error)
	GetVersionByID(ctx context.Context, id uint) (*model.Version, error)
	UpdateVersion(ctx context.Context, version *model.Version) error
	DeleteVersion(ctx context.Context, id uint) error
	CreateFeature(ctx context.Context, feature *model.Feature) error
	GetFeaturesByVersion(ctx context.Context, versionID uint) ([]model.Feature, error)
	UpdateFeature(ctx context.Context, feature *model.Feature) error
	DeleteFeature(ctx context.Context, id uint) error
	GetVersionByDeviceID(ctx context.Context, deviceID uint) (*model.Version, error)
	GetCurrnetAllPreviousVersions(ctx context.Context, deviceID uint) ([]model.Version, error)
}

type versionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) VersionRepository {
	return &versionRepository{db: db}
}

func (r *versionRepository) CreateVersion(ctx context.Context, version *model.Version) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *versionRepository) GetAllVersions(ctx context.Context) ([]model.Version, error) {
	var versions []model.Version
	err := r.db.WithContext(ctx).Preload("Features").Order("id DESC").Find(&versions).Error
	return versions, err
}

func (r *versionRepository) GetVersionByID(ctx context.Context, id uint) (*model.Version, error) {
	var version model.Version
	err := r.db.WithContext(ctx).First(&version, id).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *versionRepository) UpdateVersion(ctx context.Context, version *model.Version) error {
	return r.db.WithContext(ctx).Save(version).Error
}

func (r *versionRepository) DeleteVersion(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Version{}, id).Error
}

func (r *versionRepository) CreateFeature(ctx context.Context, feature *model.Feature) error {
	return r.db.WithContext(ctx).Create(feature).Error
}

func (r *versionRepository) GetFeaturesByVersion(ctx context.Context, versionID uint) ([]model.Feature, error) {
	var features []model.Feature
	err := r.db.WithContext(ctx).Where("version_id = ?", versionID).Find(&features).Error
	return features, err
}

func (r *versionRepository) UpdateFeature(ctx context.Context, feature *model.Feature) error {
	return r.db.WithContext(ctx).Save(feature).Error
}

func (r *versionRepository) DeleteFeature(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Feature{}, id).Error
}

func (r *versionRepository) GetVersionByDeviceID(ctx context.Context, deviceID uint) (*model.Version, error) {
	var version model.Version
	err := r.db.WithContext(ctx).
		Joins("JOIN devices ON devices.version_id = versions.id").
		Where("devices.id = ?", deviceID).
		First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *versionRepository) GetCurrnetAllPreviousVersions(ctx context.Context, deviceID uint) ([]model.Version, error) {
	var versions []model.Version
	err := r.db.WithContext(ctx).
		Joins("JOIN devices ON devices.version_id = versions.id").
		Where("devices.id = ?", deviceID).
		Or("versions.id IN (?)",
			r.db.Table("devices").Select("version_id").Where("devices.id = ?", deviceID),
		).
		Find(&versions).Error
	if err != nil {
		return nil, err
	}
	return versions, nil
}
