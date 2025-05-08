package auth

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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
	repository  user.Repository
	redisClient *redis.Client
	cfg         config.Config
}

func NewService(userRepo user.Repository, redisClient *redis.Client, cfg config.Config) Service {
	return &service{
		repository:  userRepo,
		redisClient: redisClient,
		cfg:         cfg,
	}
}

func (s *service) Login(email string, password string) (LoginResponse, error) {
	// Check if the user exists
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return LoginResponse{}, err
	}

	// Check if the password is correct
	if !utils.ValidatePassword(password, user.HashedPassword) {
		return LoginResponse{}, fmt.Errorf("invalid password")
	}

	// Convert user.ID to string
	userID := fmt.Sprintf("%d", user.ID)
	sessionId := uuid.New().String()

	// Generate access and refresh tokens
	accessToken, refreshToken, expiresAt, err := utils.GenerateJWT(s.cfg.JWT.SecretKey, sessionId, userID)
	if err != nil {
		return LoginResponse{}, err
	}

	// Session ID is the user ID
	if err != nil {
		return LoginResponse{}, err
	}

	// Store the refresh token in Redis
	err = utils.StoreTokenInRedis(s.redisClient, sessionId, accessToken)
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
	claims, err := utils.ValidateJWT(s.cfg.JWT.SecretKey, refreshToken)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	// Get the user ID and session ID from the claims
	userIDStr := fmt.Sprintf("%s", (*claims)["userID"])
	sessionId := fmt.Sprintf("%s", (*claims)["sessionId"])

	// Generate new access and refresh tokens
	accessToken, _, expiresAt, err := utils.GenerateJWT(s.cfg.JWT.SecretKey, sessionId, userIDStr)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	// Store the new refresh token in Redis
	err = utils.StoreTokenInRedis(s.redisClient, sessionId, accessToken)
	if err != nil {
		return RefreshTokenResponse{}, err
	}

	return RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *service) Logout(refreshToken string) (LogoutResponse, error) {
	// Validate the refresh token
	claims, err := utils.ValidateJWT(s.cfg.JWT.SecretKey, refreshToken)
	if err != nil {
		return LogoutResponse{}, err
	}

	// Get the session ID from the claims
	sessionId := fmt.Sprintf("%s", (*claims)["sessionId"])

	// Delete the session ID from Redis
	err = utils.DeleteTokenInRedis(s.redisClient, sessionId)

	return LogoutResponse{
		Success: true,
	}, nil
}
