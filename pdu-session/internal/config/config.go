package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Host        string
	Port        string
	LogLevel    zerolog.Level
	DatabaseURL string
}

func Load() (*Config, error) {
	godotenv.Load("../.env")
	lv := os.Getenv("LOGLEVEL")
	var loglv zerolog.Level
	switch lv {
	case "INFO":
		loglv = zerolog.InfoLevel
	case "DEBUG":
		loglv = zerolog.DebugLevel
	case "ERROR":
		loglv = zerolog.ErrorLevel
	}
	return &Config{
		Host:        os.Getenv("HOST"),
		Port:        os.Getenv("PORT"),
		LogLevel:    loglv,
		DatabaseURL: os.Getenv("DB_URL"),
	}, nil
}
