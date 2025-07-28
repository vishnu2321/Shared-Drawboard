package auth

import (
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

func CreateNewToken(username string) (string, error) {
	setEnvVariables()

	secretKey := os.Getenv("SECRETKEY_FOR_JWT")

	claims := jwt.MapClaims{
		"sub": username,
		// "role": role,
		"iss": "shared-drawboard",
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Expires in 1 hour
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
