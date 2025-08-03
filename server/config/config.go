package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Host        string
	Env         string
	DatabaseURL string
	JWTSecret   string

	// OAuth
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURL  string

	// Discord Bot
	DiscordBotToken string

	// Frontend
	FrontendURL string

	// API
	RateLimit  string
	RateWindow string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Not an error if file doesn't exist in production
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	host := getEnv("HOST", "http://localhost")
	port := getEnv("PORT", "8080")

	// Construct base API URL
	baseURL := host
	if port != "" && port != "80" && port != "443" {
		baseURL = fmt.Sprintf("%s:%s", host, port)
	}

	cfg := &Config{
		Port:        port,
		Host:        host,
		Env:         getEnv("ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),

		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:  baseURL + "/" + getEnv("GITHUB_REDIRECT_URL", "api/auth/github/callback"),

		DiscordClientID:     getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
		DiscordRedirectURL:  baseURL + "/" + getEnv("DISCORD_REDIRECT_URL", "api/auth/discord/callback"),

		DiscordBotToken: getEnv("DISCORD_BOT_TOKEN", ""),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		RateLimit:   getEnv("API_RATE_LIMIT", "100"),
		RateWindow:  getEnv("API_RATE_WINDOW", "1m"),
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
