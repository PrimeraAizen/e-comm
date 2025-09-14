package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// ErrInvalidConfig ошибка конфигурации приложения.
var ErrInvalidConfig = errors.New("invalid config")

// Путь к файлам ключей и директории миграций.
const (
	MigrationDir = "migrations"
	PathToConfig = "./config"
)

type Config struct {
	Http Http
	PG   PG
}

func LoadConfig() (*Config, error) {
	return LoadConfigFromDirectory(PathToConfig)
}

func LoadConfigFromDirectory(path string) (*Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("decode into struct: %w", err)
	}
	err = cfg.Validate()
	if err != nil {
		return nil, ErrInvalidConfig
	}

	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.Http.Host == "" {
		return fmt.Errorf("missing http host:")
	}
	if cfg.Http.Port == "" {
		return fmt.Errorf("missing http port")
	}
	return nil
}

type Http struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PG struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
	MaxConns int    `json:"max_conns"`
	MinConns int    `json:"min_conns"`
	URL      string
}
