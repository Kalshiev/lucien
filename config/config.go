package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUrl       string
	TokenSecret string
	Platform    string
	Port        string
}

// Load loads configuration from environment
func Load() (*Config, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable not set")
	}

	tokenSecret := os.Getenv("SECRET_KEY")
	if tokenSecret == "" {
		return nil, fmt.Errorf("SECRET_KEY environment variable not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "prod" // default
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	return &Config{
		DBUrl:       dbUrl,
		TokenSecret: tokenSecret,
		Platform:    platform,
		Port:        port,
	}, nil
}
