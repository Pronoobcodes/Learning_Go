package config

import "os"

const (
	defaultPort        = "8080"
	defaultDatabaseURL = "postgres://app:password@localhost:5433/url_shortener?sslmode=disable"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}, nil
}
