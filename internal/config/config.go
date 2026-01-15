package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	DatabaseURL string `yaml:"database_url"`
	LLMURL      string `yaml:"llm_url"`
	LLMKey      string `yaml:"llm_key"`
	LLMModel    string `yaml:"llm_model"`
	LogPath     string `yaml:"log_path"`
	LogLevel    string `yaml:"log_level"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		DatabaseURL: "postgres://user:password@192.168.64.2:5432/bd_bot?sslmode=disable",
		LLMURL:      "http://127.0.0.1:11434/v1",
		LLMKey:      "sk-...",
		LLMModel:    "gemma3:4b",
		LogPath:     "bd_bot.log",
		LogLevel:    "INFO",
	}
}

// LoadConfig reads configuration from config.yaml or returns default
func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile("config.yaml")
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}