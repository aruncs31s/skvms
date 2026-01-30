package service

import (
	"context"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	List(ctx context.Context) ([]dto.UserView, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) error
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
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
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
	}
	return views, nil
}

func (s *userService) Create(ctx context.Context, req *dto.CreateUserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) error {
	user := &model.User{
		ID:       id,
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
