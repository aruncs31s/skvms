package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type AuditService interface {
	Log(ctx context.Context, userID uint, username, action, details, ipAddress string) error
	List(ctx context.Context, action string, limit int) ([]model.AuditLog, error)
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Log(ctx context.Context, userID uint, username, action, details, ipAddress string) error {
	log := &model.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    action,
		Details:   details,
		IPAddress: ipAddress,
	}
	return s.repo.Create(ctx, log)
}

func (s *auditService) List(ctx context.Context, action string, limit int) ([]model.AuditLog, error) {
	return s.repo.List(ctx, action, limit)
}
