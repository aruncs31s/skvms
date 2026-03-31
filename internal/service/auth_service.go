package service

import (
	"context"
	"errors"
	"time"

	"github.com/aruncs31s/skvms/internal/dto"
	"github.com/aruncs31s/skvms/internal/model"
	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, string, *model.User, error)
	Register(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *authService) Login(
	ctx context.Context,
	username, password string) (string, string, *model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", nil, err
	}
	if user == nil {
		return "", "", nil, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, nil
	}

	return s.generateTokenPair(user)
}

func (s *authService) generateTokenPair(user *model.User) (string, string, *model.User, error) {
	// Access Token
	accessClaims := jwt.MapClaims{
		"sub":        user.ID,
		"username":   user.Username,
		"token_type": "access",
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
		"iat":        time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"sub":        user.ID,
		"username":   user.Username,
		"token_type": "refresh",
		"exp":        time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	return accessTokenString, refreshTokenString, user, nil
}

func (s *authService) Refresh(ctx context.Context, refreshTokenString string) (string, string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	if tokenType, ok := claims["token_type"].(string); !ok || tokenType != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", "", errors.New("invalid token format")
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	accessToken, newRefreshToken, _, err := s.generateTokenPair(user)
	return accessToken, newRefreshToken, err
}

func (s *authService) Register(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
