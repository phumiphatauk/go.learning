package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// Middleware to validate JWT token
// Middleware to validate JWT token
func TokenAuthMiddleware(next echo.HandlerFunc, redisClient *redis.Client, jwt_secret_key string) echo.HandlerFunc {
	var jwtSecretKey = []byte(jwt_secret_key)
	return func(c echo.Context) error {
		// Get the token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is missing")
		}

		// Bearer token format: "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := tokenParts[1]

		// Validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		// Extract claims for checking session ID
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Retrieve sessionID from claims
			sessionID, sessionExists := claims["sessionId"].(string)
			if !sessionExists || sessionID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Session ID is missing from token")
			}

			// Check sessionID in Redis
			storedToken, err := redisClient.Get(c.Request().Context(), sessionID).Result()
			if err == redis.Nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Session ID is invalid or expired")
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to validate session ID")
			}

			// Ensure the token matches the stored token in Redis
			if storedToken != tokenString {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token does not match session")
			}

			// Set userID in the context for later use
			c.Set("userID", claims["userID"])
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "Failed to parse token claims")
		}

		// Continue to the next handler
		return next(c)
	}
}
