package config

import (
	"errors"
	"fmt"
	"os"
)

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
	DBURL      string
	RedisAddr  string
	HashConfig hashConfig
	IsDev      bool
}

func LoadConfig() (*config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		fmt.Println("Environment Variable 'SERVER_PORT' not found, using default: 8080")
		serverPort = "8080"
	}

	// "postgres://postgres:postgres@localhost:5433/db"
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbPort == "" || dbName == "" {
		return nil, errors.New("failed to get env")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		return nil, errors.New("failed to get REDIS_PORT env")
	}

	dev := os.Getenv("DEV")

	return &config{
		ServerPort: serverPort,
		DBURL: fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			dbPort,
			dbName,
		),
		RedisAddr: fmt.Sprint("localhost:", redisPort),
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
