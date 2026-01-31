package service

import (
	"context"
	"time"

	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *model.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, *model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, nil
	}

	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
