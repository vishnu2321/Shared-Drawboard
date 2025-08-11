package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/shared-drawboard/pkg/logger"
)

func setEnvVariables() {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading environment variables: %s", err)
		os.Exit(1)
	}
}

func CreateJWTToken(username string, expiryTime int64) (string, error) {
	setEnvVariables()

	secretKey := os.Getenv("SECRETKEY_FOR_JWT")

	claims := jwt.MapClaims{
		"sub": username,
		// "role": role,
		"iss": "shared-drawboard",
		"exp": expiryTime,
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWTToken(tokenStr string) (jwt.Claims, error) {
	setEnvVariables()
	secretKey := os.Getenv("SECRETKEY_FOR_JWT")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return jwt.MapClaims{}, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
