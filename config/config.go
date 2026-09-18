// Package config loads the application configuration from environment
// variables (SERVER_PORT, DB_USER/DB_PASSWORD/DB_PORT/DB_NAME, REDIS_PORT,
// DEV) and exposes it to the rest of the application.
package config

import (
	"errors"
	"fmt"
	"os"
)

var serverPort = "8080"

const (
	time    uint32 = 1
	memory  uint32 = 64 * 1024
	threads uint8  = 4
	keyLen  uint32 = 32
	saltLen uint32 = 16
)

type hashConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

type config struct {
	ServerPort string
	// postgres://<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>
	DBURL      string
	RedisAddr  string
	HashConfig hashConfig
	IsDev      bool
}

// LoadConfig reads the configuration from environment variables and returns
// a populated config, or an error when a required variable is missing.
func LoadConfig() (*config, error) {
	dev := os.Getenv("DEV")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, errors.New("failed to get postgres environments")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		return nil, errors.New("failed to get redis environments")
	}

	return &config{
		ServerPort: serverPort,
		DBURL: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		),
		RedisAddr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		HashConfig: hashConfig{
			Time:    time,
			Memory:  memory,
			Threads: threads,
			KeyLen:  keyLen,
			SaltLen: saltLen,
		},
		IsDev: dev == "true",
	}, nil
}
