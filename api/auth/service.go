package auth

import (
	"fmt"

	"go.learning/api/user"
	"go.learning/config"
	"go.learning/utils"
)

type Service interface {
	Login(email, password string) (LoginResponse, error)
	RefreshToken(refreshToken string) (RefreshTokenResponse, error)
	Logout(refreshToken string) (LogoutResponse, error)
}
type service struct {
	user.Repository
	cfg config.Config
}

func NewService(userRepo user.Repository, cfg config.Config) Service {
	return &service{
		Repository: userRepo,
		cfg:        cfg,
	}
}

func (s *service) Login(email string, password string) (LoginResponse, error) {
	// Check if the user exists
	user, err := s.Repository.GetUserByEmail(email)
	if err != nil {
		return LoginResponse{}, err
	}

	// Check if the password is correct
	if !utils.ValidatePassword(password, user.HashedPassword) {
		return LoginResponse{}, fmt.Errorf("invalid password")
	}

	// Convert user.ID to string
	userID := fmt.Sprintf("%d", user.ID)

	// Generate access and refresh tokens
	accessToken, refreshToken, expiresAt, err := utils.GenerateJWT(s.cfg.JWT.SecretKey, userID)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *service) RefreshToken(refreshToken string) (RefreshTokenResponse, error) {
	// Validate the refresh token
	userID, err := utils.ValidateJWT(s.cfg.JWT.SecretKey, refreshToken)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	// Convert userID to string
	userIDStr := fmt.Sprintf("%d", userID)

	// Generate new access and refresh tokens
	accessToken, newRefreshToken, expiresAt, err := utils.GenerateJWT(s.cfg.JWT.SecretKey, userIDStr)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	return RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *service) Logout(refreshToken string) (LogoutResponse, error) {
	// Validate the refresh token
	_, err := utils.ValidateJWT(refreshToken, s.cfg.JWT.SecretKey)
	if err != nil {
		return LogoutResponse{}, err
	}

	// TODO: Invalidate the refresh token in the database or cache

	return LogoutResponse{
		Success: true,
	}, nil
}
