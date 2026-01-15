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
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		DatabaseURL: "postgres://user:password@localhost:5432/bd_bot?sslmode=disable",
		LLMURL:      "http://localhost:8000/v1",
		LLMKey:      "sk-...",
		LLMModel:    "gemma3:4b",
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