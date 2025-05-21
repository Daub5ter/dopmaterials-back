package config

import (
	"os"
	"time"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

// Config - конфиг приложения.
type Config struct {
	Server   serverConfig   `yaml:"server"`
	Database databaseConfig `yaml:"database"`
	Search   searchConfig   `yaml:"search"`
	Logger   loggerConfig   `yaml:"logger"`
}

// serverConfig - конфиг сервера.
type serverConfig struct {
	Domain  string        `yaml:"domain"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// databaseConfig - конфиг базы данных.
type databaseConfig struct {
	Host     string        `env:"POSTGRES_HOST"`
	Port     int           `env:"POSTGRES_PORT"`
	User     string        `env:"POSTGRES_USER"`
	Password string        `env:"POSTGRES_PASSWORD"`
	DBName   string        `env:"POSTGRES_DB_NAME"`
	SSLMode  string        `env:"POSTGRES_SSL_MODE"`
	Timeout  time.Duration `yaml:"timeout"`
}

// searchConfig - конфиг поисковой системы
type searchConfig struct {
	Username  string        `env:"SEARCH_USERNAME"`
	Password  string        `env:"SEARCH_PASSWORD"`
	Index     string        `yaml:"index"`
	Addresses []string      `yaml:"addresses"`
	Timeout   time.Duration `yaml:"timeout"`
}

// loggerConfig - структура конфиграции логов.
type loggerConfig struct {
	Level string `yaml:"level"`
}

// NewConfig создает API конфига.
func NewConfig(configPath string) (*Config, error) {
	// Считывание файла конфигурации.
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Декодирование данных в структуру конфигурации.
	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&cfg.Database)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&cfg.Search)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
