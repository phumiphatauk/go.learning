package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func StoreTokenInRedis(redisClient *redis.Client, sessionID, token string) error {
	// Set the expiration time for the token
	err := redisClient.Set(context.Background(), sessionID, token, time.Hour*24).Err()
	return err
}

func DeleteTokenInRedis(redisClient *redis.Client, sessionID string) error {
	// Delete the token from Redis
	err := redisClient.Del(context.Background(), sessionID).Err()
	return err
}

func CheckSessionInRedis(redisClient *redis.Client, sessionID string) (string, error) {
	// Check if the session ID exists in Redis
	return redisClient.Get(context.Background(), sessionID).Result()
}
