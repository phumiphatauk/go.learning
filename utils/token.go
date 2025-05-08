package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT generates a JWT token
func GenerateJWT(secretKey string, sessionId string, userID string) (string, string, int64, error) {

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"sessionId": sessionId,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", 0, err
	}

	// Create a refresh token 7 days from now
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"sessionId": sessionId,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", 0, err
	}

	// Create the expiration time
	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	if err != nil {
		return "", "", 0, err
	}

	return signedToken, refreshTokenString, expiresAt, nil
}

func ValidateJWT(secretKey string, tokenString string) (*jwt.MapClaims, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	// Handle validation errors
	if err != nil {
		return nil, err
	}

	// Extract and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
