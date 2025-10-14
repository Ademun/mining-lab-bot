package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PollingConfig  PollingConfig  `yaml:"polling"`
	TelegramConfig TelegramConfig `yaml:"telegram"`
}

type PollingConfig struct {
	ServiceURL string      `yaml:"service_url"`
	Mode       PollingMode `yaml:"mode"`
}

type PollingMode string

const (
	ModeNormal     PollingMode = "normal"
	ModeAggressive PollingMode = "aggressive"
)

type TelegramConfig struct {
	BotToken string
	AdminID  int
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) loadFromEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	if botToken := os.Getenv("BOT_TOKEN"); botToken != "" {
		c.TelegramConfig.BotToken = botToken
	}

	if adminID, err := strconv.Atoi(os.Getenv("ADMIN_ID")); err == nil {
		c.TelegramConfig.AdminID = adminID
	}

	return nil
}
