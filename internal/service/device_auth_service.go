package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aruncs31s/skvms/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type DeviceAuthService interface {
	GenerateDeviceToken(ctx context.Context, userID uint, deviceID uint) (string, error)
	ValidateDeviceToken(ctx context.Context, tokenString string) (*DeviceTokenClaims, error)
}

type deviceAuthService struct {
	deviceRepo repository.DeviceRepository
	userRepo   repository.UserRepository
	jwtSecret  []byte
}

type DeviceTokenClaims struct {
	UserID   uint `json:"user_id"`
	DeviceID uint `json:"device_id"`
	jwt.RegisteredClaims
}

func NewDeviceAuthService(deviceRepo repository.DeviceRepository, userRepo repository.UserRepository, jwtSecret string) DeviceAuthService {
	return &deviceAuthService{
		deviceRepo: deviceRepo,
		userRepo:   userRepo,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (s *deviceAuthService) GenerateDeviceToken(ctx context.Context, userID uint, deviceID uint) (string, error) {
	// Verify that the user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to verify user: %w", err)
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Verify that the device exists
	device, err := s.deviceRepo.GetDevice(ctx, deviceID)
	if err != nil {
		return "", fmt.Errorf("failed to verify device: %w", err)
	}
	if device == nil {
		return "", errors.New("device not found")
	}

	// Verify that the device belongs to the user
	if device.CreatedBy != userID {
		return "", errors.New("device does not belong to user")
	}

	// Create JWT claims with UserID and DeviceID
	claims := DeviceTokenClaims{
		UserID:   userID,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2400 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("device:%d:user:%d", deviceID, userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *deviceAuthService) ValidateDeviceToken(ctx context.Context, tokenString string) (*DeviceTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DeviceTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*DeviceTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
