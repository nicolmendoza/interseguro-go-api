package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string
	JWTSecret  string
	NodeAPIURL string
}

func Load() (Config, error) {

	port, err := requiredEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	nodeAPIURL, err := requiredEnv("NODE_API_URL")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:       port,
		JWTSecret:  jwtSecret,
		NodeAPIURL: nodeAPIURL,
	}, nil
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("falta la variable de entorno requerida: %s", key)
	}
	return value, nil
}
