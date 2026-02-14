package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
)

type AuditService interface {
	Log(ctx context.Context, userID uint, username, action, details, ipAddress string) error
	LogDeviceAction(
		ctx context.Context,
		userID uint,
		username,
		action,
		details,
		ipAddress string,
		deviceID uint,
	) error
	List(
		ctx context.Context,
		action string,
		limit int,
	) ([]model.AuditLog, error)
	ListByUser(ctx context.Context, userID uint, limit int) ([]model.AuditLog, error)
}

type auditService struct {
	repo repository.AuditRepository
}

func NewAuditService(repo repository.AuditRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Log(
	ctx context.Context,
	userID uint,
	username,
	action,
	details,
	ipAddress string,
) error {
	// Find action ID from string
	var actionID model.DeviceAction
	for id, name := range model.DeviceActionsMap {
		if name == action {
			actionID = id
			break
		}
	}
	log := &model.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    actionID,
		Details:   details,
		IPAddress: ipAddress,
	}
	return s.repo.Create(ctx, log)
}
func (s *auditService) LogDeviceAction(
	ctx context.Context,
	userID uint,
	username,
	action,
	details,
	ipAddress string,
	deviceID uint,
) error {
	// Find action ID from string
	var actionID model.DeviceAction
	for id, name := range model.DeviceActionsMap {
		if name == action {
			actionID = id
			break
		}
	}
	log := &model.AuditLog{
		UserID:    userID,
		Username:  username,
		Action:    actionID,
		Details:   details,
		IPAddress: ipAddress,
		DeviceID:  &deviceID,
	}
	return s.repo.Create(ctx, log)
}
func (s *auditService) List(ctx context.Context, action string, limit int) ([]model.AuditLog, error) {
	return s.repo.List(ctx, action, limit)
}

func (s *auditService) ListByUser(ctx context.Context, userID uint, limit int) ([]model.AuditLog, error) {
	return s.repo.ListByUser(ctx, userID, limit)
}
