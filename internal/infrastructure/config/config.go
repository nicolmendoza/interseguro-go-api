package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port                   string
	JWTSecret              string
	FrontendURL            string
	NodeAPIURL             string
	RateLimitMax           int
	RateLimitWindowSeconds int
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

	frontendURL, err := requiredEnv("FRONTEND_URL")
	if err != nil {
		return Config{}, err
	}

	nodeAPIURL, err := requiredEnv("NODE_API_URL")
	if err != nil {
		return Config{}, err
	}

	rateLimitMax, err := requiredIntEnv("RATE_LIMIT_MAX")
	if err != nil {
		return Config{}, err
	}

	rateLimitWindowSeconds, err := requiredIntEnv("RATE_LIMIT_WINDOW_SECONDS")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:                   port,
		JWTSecret:              jwtSecret,
		FrontendURL:            frontendURL,
		NodeAPIURL:             nodeAPIURL,
		RateLimitMax:           rateLimitMax,
		RateLimitWindowSeconds: rateLimitWindowSeconds,
	}, nil
}

func requiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("falta la variable de entorno requerida: %s", key)
	}
	return value, nil
}

func requiredIntEnv(key string) (int, error) {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("la variable de entorno %s debe ser un numero positivo", key)
	}
	return value, nil
}
