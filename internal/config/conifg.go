// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GeminiAPIKey  string    `yaml:"gemini_api_key"`
	DefaultModels []string  `yaml:"default_models"`
	MaxTokens     int       `yaml:"max_tokens"`
	RateLimit     RateLimit `yaml:"rate_limit"`
}

type RateLimit struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Read config from YAML
	configFile, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Error parsing config:", err)
	}

	// Override with environment variable if set
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		config.GeminiAPIKey = apiKey
	}

	// Default models if not specified
	if len(config.DefaultModels) == 0 {
		config.DefaultModels = []string{
			"gemini-1.0-pro",
			"gemini-pro",
			"gemini-1.5-pro",
		}
	}

	// Default rate limiting if not specified
	if config.RateLimit.RequestsPerMinute == 0 {
		config.RateLimit.RequestsPerMinute = 5 // Default: 5 requests per minute
	}

	return &config
}
