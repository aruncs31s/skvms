package service

import (
	"context"
	"errors"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	List(ctx context.Context) ([]dto.UserView, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) error
	Delete(ctx context.Context, id uint) error
	UserReader
	GetProfile(ctx context.Context, userID uint) (*dto.UserProfile, error)
}
type UserReader interface {
	GetByID(ctx context.Context, id uint) (*dto.UserView, error)
}

type userService struct {
	repo          repository.UserRepository
	deviceService DeviceService
	auditService  AuditService
}

func NewUserService(repo repository.UserRepository, deviceService DeviceService, auditService AuditService) UserService {
	return &userService{
		repo:          repo,
		deviceService: deviceService,
		auditService:  auditService,
	}
}

func (s *userService) List(ctx context.Context) ([]dto.UserView, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	views := make([]dto.UserView, len(users))
	for i, user := range users {
		views[i] = dto.UserView{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
	}
	return views, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*dto.UserView, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	view := &dto.UserView{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
	return view, nil
}

func (s *userService) Create(ctx context.Context, req *dto.CreateUserRequest) error {

	if req.Username == "" || req.Password == "" {
		return errors.New("username and password are required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	role := req.Role
	if role == "" {
		role = "user" // default role
	}
	name := req.Name
	if name == "" {
		name = "User_" + req.Username
	}
	user := &model.User{
		Name:     name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		logger.GetLogger().Error("Failed to create user: ", zap.String("username", req.Username), zap.Error(err))
	}
	return err
}

func (s *userService) Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) error {
	user := &model.User{
		ID:       id,
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*dto.UserProfile, error) {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	devices, err := s.deviceService.ListDevicesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	auditLogs, err := s.auditService.ListByUser(ctx, userID, 50) // limit to 50 recent activities
	if err != nil {
		return nil, err
	}

	activity := make([]dto.AuditLogView, len(auditLogs))
	for i, log := range auditLogs {
		actionName := "unknown"
		if name, exists := model.DeviceActionsMap[log.Action]; exists {
			actionName = name
		}
		activity[i] = dto.AuditLogView{
			ID:        log.ID,
			Username:  log.Username,
			Action:    actionName,
			Details:   log.Details,
			IPAddress: log.IPAddress,
			DeviceID:  log.DeviceID,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &dto.UserProfile{
		User:     *user,
		Devices:  devices,
		Activity: activity,
	}, nil
}
