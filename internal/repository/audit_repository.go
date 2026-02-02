package repository

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"gorm.io/gorm"
)

type AuditRepository interface {
	Create(
		ctx context.Context,
		log *model.AuditLog,
	) error
	List(ctx context.Context, action string, limit int) ([]model.AuditLog, error)
	ListByUser(ctx context.Context, userID uint, limit int) ([]model.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditRepository) List(ctx context.Context, action string, limit int) ([]model.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	query := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit)
	if action != "" {
		query = query.Where("action = ?", action)
	}

	var logs []model.AuditLog
	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditRepository) ListByUser(ctx context.Context, userID uint, limit int) ([]model.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var logs []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
