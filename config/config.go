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
	Http Http `mapstructure:"http"`
	PG   PG   `mapstructure:"database"`
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

	cfg.PG.URL = cfg.PG.connString()
	err = cfg.Validate()
	if err != nil {
		return nil, err
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
	if cfg.PG.URL == "" {
		return fmt.Errorf("missing database url")
	}
	return nil
}

func (d *PG) connString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.SSLMode)
}

type Http struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
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
