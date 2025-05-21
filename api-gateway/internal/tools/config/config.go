package config

import (
	"os"
	"time"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

// Config - конфиг приложения.
type Config struct {
	Server  serverConfig  `yaml:"server"`
	Logger  loggerConfig  `yaml:"logger"`
	GRPC    gRPCConfig    `yaml:"grpc"`
	Limiter limiterConfig `yaml:"limiter"`
}

// serverConfig - конфиг сервера.
type serverConfig struct {
	Domain  string        `yaml:"domain"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type service struct {
	Connection string        `yaml:"connection"`
	Timeout    time.Duration `yaml:"timeout"`
}

type gRPCConfig struct {
	Content service `yaml:"content"`
}

type limiterConfig struct {
	MaxRequestsTimeLivingIpAddress int           `yaml:"max_requests_time_living_ip_address"`
	Addr                           string        `yaml:"addr"`
	Password                       string        `env:"REDIS_PASSWORD"`
	Timeout                        time.Duration `yaml:"timeout"`
	TimeLivingIpAddress            time.Duration `yaml:"time_living_ip_address"`
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

	err = env.Parse(&cfg.Limiter)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
