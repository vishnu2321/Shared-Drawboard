package database

import (
	"os"

	dotenv "github.com/joho/godotenv"
	"github.com/shared-drawboard/pkg/logger"
)

func Config() map[string]interface{} {
	err := dotenv.Load()
	if err != nil {
		logger.Error("DB Config: Error loading environment variables")
	}
	db_uri := os.Getenv("MONGODB_URI")
	db_name := os.Getenv("MONGODB_NAME")
	return map[string]interface{}{"uri": db_uri, "db_name": db_name}
}
