package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PollingConfig  PollingConfig  `yaml:"polling"`
	TelegramConfig TelegramConfig `yaml:"telegram"`
	TeacherConfig  TeacherConfig  `yaml:"teacher"`
}

type PollingConfig struct {
	ServiceURL          string        `yaml:"service_url"`
	Mode                PollingMode   `yaml:"mode"`
	ServiceIDUpdateRate time.Duration `yaml:"service_id_update_rate"`
	NormalPollRate      time.Duration `yaml:"normal_poll_rate"`
	AggressivePollRate  time.Duration `yaml:"aggressive_poll_rate"`
	NormalFetchRate     time.Duration `yaml:"normal_fetch_rate"`
	AggressiveFetchRate time.Duration `yaml:"aggressive_fetch_rate"`
	MinFetchRate        time.Duration `yaml:"min_fetch_rate"`
	MaxFetchRate        time.Duration `yaml:"max_fetch_rate"`
	BackoffFactor       float64       `yaml:"backoff_factor"`
	RecoveryFactor      float64       `yaml:"recovery_factor"`
}

func (s *PollingConfig) GetFetchRate() time.Duration {
	if s.Mode == ModeAggressive {
		return s.AggressiveFetchRate
	}
	return s.NormalFetchRate
}

func (s *PollingConfig) SetFetchRate(rate time.Duration) {
	if s.Mode == ModeAggressive {
		s.AggressiveFetchRate = rate
	}
	s.NormalFetchRate = rate
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

type TeacherConfig struct {
	StartingWeek int `yaml:"starting_week"`
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
