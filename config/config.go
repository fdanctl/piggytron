package config

import (
	"errors"
	"fmt"
	"os"
)

type config struct {
	ServerPort string
	DBURL      string
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

	return &config{
		ServerPort: serverPort,
		DBURL: fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s",
			dbUser,
			dbPassword,
			dbPort,
			dbName,
		),
	}, nil
}
