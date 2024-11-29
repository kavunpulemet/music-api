package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

const (
	serverPortEnv     = "SERVER_PORT"
	dbHostEnv         = "DB_HOST"
	dbPortEnv         = "DB_PORT"
	dbUserEnv         = "DB_USER"
	dbNameEnv         = "DB_NAME"
	dbPasswordEnv     = "DB_PASSWORD"
	dbSSLModeEnv      = "DB_SSLMODE"
	externalAPIUrlEnv = "EXTERNAL_API_URL"
)

type Config struct {
	ServerPort         string
	DBConnectionString string
	SongDetailsAPIUrl  string
	LoggerLevel        zapcore.Level
}

func NewConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	loggerLevelStr := os.Getenv("LOGGER_LEVEL")
	var loggerLevel zapcore.Level
	err := loggerLevel.UnmarshalText([]byte(loggerLevelStr))
	if err != nil {
		fmt.Printf("invalid log level: %s. Defaulting to info\n", loggerLevelStr)
		loggerLevel = zapcore.InfoLevel
	}

	return Config{
		ServerPort: fmt.Sprintf(":%s", os.Getenv(serverPortEnv)),
		DBConnectionString: fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			os.Getenv(dbHostEnv), os.Getenv(dbPortEnv), os.Getenv(dbUserEnv),
			os.Getenv(dbNameEnv), os.Getenv(dbPasswordEnv), os.Getenv(dbSSLModeEnv)),
		SongDetailsAPIUrl: os.Getenv(externalAPIUrlEnv),
		LoggerLevel:       loggerLevel,
	}, nil
}
