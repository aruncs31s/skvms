package service

import (
	"context"
	"errors"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type VersionService interface {
	CreateVersion(ctx context.Context, version string) (*model.Version, error)
	GetAllVersions(ctx context.Context) ([]dto.VersionResponse, error)
	GetVersionByID(ctx context.Context, id uint) (*model.Version, error)
	UpdateVersion(ctx context.Context, id uint, version string) (*model.Version, error)
	DeleteVersion(ctx context.Context, id uint) error
	CreateFeature(ctx context.Context, versionID uint, featureName string, enabled bool) (*model.Feature, error)
	GetFeaturesByVersion(ctx context.Context, versionID uint) ([]model.Feature, error)
	UpdateFeature(ctx context.Context, id uint, featureName string, enabled bool) (*model.Feature, error)
	DeleteFeature(ctx context.Context, id uint) error
}

type versionService struct {
	repo repository.VersionRepository
}

func NewVersionService(repo repository.VersionRepository) VersionService {
	return &versionService{repo: repo}
}

func (s *versionService) CreateVersion(ctx context.Context, version string) (*model.Version, error) {
	if version == "" {
		return nil, errors.New("version cannot be empty")
	}

	v := &model.Version{Version: version}
	err := s.repo.CreateVersion(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *versionService) GetAllVersions(ctx context.Context) ([]dto.VersionResponse, error) {
	versions, err := s.repo.GetAllVersions(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.VersionResponse
	for _, v := range versions {
		var features []dto.FeatureResponse
		for _, f := range v.Features {
			features = append(features, dto.FeatureResponse{
				ID:          f.ID,
				VersionID:   f.VersionID,
				FeatureName: f.FeatureName,
				Enabled:     f.Enabled,
			})
		}
		responses = append(responses, dto.VersionResponse{
			ID:        v.ID,
			Version:   v.Version,
			CreatedAt: v.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Features:  features,
		})
	}
	return responses, nil
}

func (s *versionService) GetVersionByID(ctx context.Context, id uint) (*model.Version, error) {
	return s.repo.GetVersionByID(ctx, id)
}

func (s *versionService) UpdateVersion(ctx context.Context, id uint, version string) (*model.Version, error) {
	if version == "" {
		return nil, errors.New("version cannot be empty")
	}

	v, err := s.repo.GetVersionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	v.Version = version
	err = s.repo.UpdateVersion(ctx, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s *versionService) DeleteVersion(ctx context.Context, id uint) error {
	return s.repo.DeleteVersion(ctx, id)
}

func (s *versionService) CreateFeature(ctx context.Context, versionID uint, featureName string, enabled bool) (*model.Feature, error) {
	if featureName == "" {
		return nil, errors.New("feature name cannot be empty")
	}

	// Check if version exists
	_, err := s.repo.GetVersionByID(ctx, versionID)
	if err != nil {
		return nil, err
	}

	feature := &model.Feature{
		VersionID:   versionID,
		FeatureName: featureName,
		Enabled:     enabled,
	}
	err = s.repo.CreateFeature(ctx, feature)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

func (s *versionService) GetFeaturesByVersion(ctx context.Context, versionID uint) ([]model.Feature, error) {
	return s.repo.GetFeaturesByVersion(ctx, versionID)
}

func (s *versionService) UpdateFeature(ctx context.Context, id uint, featureName string, enabled bool) (*model.Feature, error) {
	if featureName == "" {
		return nil, errors.New("feature name cannot be empty")
	}

	feature := &model.Feature{
		ID:          id,
		FeatureName: featureName,
		Enabled:     enabled,
	}
	err := s.repo.UpdateFeature(ctx, feature)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

func (s *versionService) DeleteFeature(ctx context.Context, id uint) error {
	return s.repo.DeleteFeature(ctx, id)
}
