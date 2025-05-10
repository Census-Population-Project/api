package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DevMode bool
	ApiPath string
	Server  struct {
		Host         string
		Port         int
		AllowOrigins []string
	}
	Database struct {
		Name           string
		Username       string
		Password       string
		Host           string
		Port           int
		SSLMode        bool
		MaxConnections int
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
	Version string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	err := LoadEnv()
	if err != nil {
		return cfg, err
	}

	cfg.DevMode = os.Getenv("DEV_MODE") == "true"

	cfg.ApiPath = os.Getenv("API_PATH")

	cfg.Server.Host = "0.0.0.0"
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	cfg.Server.Port = 3000
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port, _ = strconv.Atoi(port)
	}

	cfg.Server.AllowOrigins = []string{"*"}
	if origins := os.Getenv("HTTP_ALLOW_ORIGINS"); origins != "" {
		cfg.Server.AllowOrigins = strings.Split(origins, ",")
	}

	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.Username = os.Getenv("DB_USERNAME")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Printf("Error parsing DB_PORT: %v", err)
		return nil, err
	}
	cfg.Database.SSLMode = os.Getenv("DB_SSL_MODE") == "true"
	maxConnsStr := os.Getenv("DB_MAX_CONNECTIONS")
	var psqlMaxConns int
	if maxConnsStr == "" {
		psqlMaxConns = 10
	} else {
		var err error
		psqlMaxConns, err = strconv.Atoi(maxConnsStr)
		if err != nil {
			log.Printf("Error parsing DB_MAX_CONNECTIONS: %v", err)
			return nil, err
		}
	}
	cfg.Database.MaxConnections = psqlMaxConns

	cfg.Redis.DB, err = strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Printf("Error parsing REDIS_AUTH_DB: %v", err)
		return nil, err
	}
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Printf("Error parsing REDIS_PORT: %v", err)
		return nil, err
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	return cfg, nil
}

func LoadEnv() error {
	//err := godotenv.Load()

	if os.Getenv("DEV_MODE") == "" {
		return fmt.Errorf("environment variables is not set")
	}

	return nil
}
