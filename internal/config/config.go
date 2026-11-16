package config

import "os"

type Config struct {
	APIAddr      string
	DatabasePath string
	CORSOrigins  string
}

func Load() Config {
	return Config{
		APIAddr:      env("API_ADDR", ":8080"),
		DatabasePath: env("DATABASE_PATH", "studybuddy.db"),
		CORSOrigins:  env("CORS_ORIGINS", "http://localhost:5173,http://localhost:4173,http://127.0.0.1:5173"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
