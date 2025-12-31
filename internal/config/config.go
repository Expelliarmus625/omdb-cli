package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	APIKey     string `json:"api_key"`
	APIBaseUrl string `json:"api_baseurl"`
	APITimeout string `json:"api_timeout"` // API timeout in seconds
}

func Load(path string) (*Config, error) {
	if apiKey := os.Getenv("OMDB_API_KEY"); apiKey != "" {
		return &Config{
			APIKey:     apiKey,
			APIBaseUrl: getEnvOrDefault("OMDB_BASE_URL", "https://www.omdbapi.com/"),
			APITimeout: getEnvOrDefault("OMDB_API_TIMEOUT", "5"),
		}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
