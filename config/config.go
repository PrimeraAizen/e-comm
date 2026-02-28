package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/PrimeraAizen/e-comm/pkg/logger"
)

// ErrInvalidConfig ошибка конфигурации приложения.
var ErrInvalidConfig = errors.New("invalid config")

// Путь к файлам ключей и директории миграций.
const (
	MigrationDir = "migrations"
	PathToConfig = "./config"
)

type Config struct {
	Http   Http          `mapstructure:"http"`
	Logger logger.Config `mapstructure:"logger"`
	JWT    JWT           `mapstructure:"jwt"`
	PG     PG            `mapstructure:"database"`
}

func (d *PG) connString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.SSLMode)
}

func LoadConfig() (*Config, error) {
	return LoadConfigFromDirectory(PathToConfig)
}

func LoadConfigFromDirectory(path string) (*Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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

	cfg.PG.URL = cfg.PG.connString()

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.Http.Host == "" {
		return fmt.Errorf("%w: missing http host", ErrInvalidConfig)
	}
	if cfg.Http.Port == "" {
		return fmt.Errorf("%w: missing http port", ErrInvalidConfig)
	}
	if cfg.PG.URL == "" {
		return fmt.Errorf("%w: missing database url", ErrInvalidConfig)
	}
	// Set default logger config if not provided
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = logger.LevelInfo
	}
	if cfg.Logger.Format == "" {
		cfg.Logger.Format = "json"
	}
	if cfg.Logger.Output == "" {
		cfg.Logger.Output = "stdout"
	}
	if cfg.Logger.Service == "" {
		cfg.Logger.Service = "template"
	}
	if cfg.Logger.Version == "" {
		cfg.Logger.Version = "1.0.0"
	}
	if cfg.Logger.Environment == "" {
		cfg.Logger.Environment = "development"
	}

	// JWT config validation
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("missing jwt secret")
	}
	if cfg.JWT.AccessTokenDuration == "" {
		cfg.JWT.AccessTokenDuration = "15m"
	}
	if cfg.JWT.RefreshTokenDuration == "" {
		cfg.JWT.RefreshTokenDuration = "168h"
	}

	return nil
}

type Http struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type JWT struct {
	Secret               string `mapstructure:"secret"`
	AccessTokenDuration  string `mapstructure:"access_token_duration"`
	RefreshTokenDuration string `mapstructure:"refresh_token_duration"`
}

type PG struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConns int    `mapstructure:"max_conns"`
	MinConns int    `mapstructure:"min_conns"`
	URL      string
}
