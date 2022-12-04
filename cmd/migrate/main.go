package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"stockinos.com/api/models"
	"stockinos.com/api/storage"
	"stockinos.com/api/utils"
)

func main() {
	os.Exit(start())
}

func start() int {
	log, err := createLogger("development")
	if err != nil {
		log.Error("Error setting up the logger", zap.Error(err))
		return 1
	}

	db := storage.NewDatabase(storage.NewDatabaseOptions{
		Host:                  utils.GetDefault("DB_HOST", "localhost"),
		Port:                  utils.GetIntDefault("DB_PORT", 5432),
		User:                  utils.GetDefault("DB_USER", "stockinos"),
		Password:              utils.GetDefault("DB_PASSWORD", "stockinos"),
		Name:                  utils.GetDefault("DB_NAME", "stockinos"),
		MaxOpenConnections:    utils.GetIntDefault("DB_MAX_OPEN_CONNECTION", 10),
		MaxIdleConnections:    utils.GetIntDefault("DB_MAX_IDLE_CONNECTION", 10),
		ConnectionMaxLifetime: utils.GetDurationDefault("DB_CONNECTION_MAX_LIFETIME", time.Hour),
		Log:                   zap.NewNop(),
	})
	if err := db.Connect(); err != nil {
		log.Error("Error connection to database", zap.Error(err))
		return 1
	}

	db.DB.AutoMigrate(
		&models.WhatsAppMessage{},
		&models.WhatsAppMessageText{},
		&models.WhatsAppMessageImage{},
		&models.WhatsAppMessageAudio{},
	)

	return 0
}

func createLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}
