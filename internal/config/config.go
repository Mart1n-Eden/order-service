package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Nats   NatsConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	DBName   string
	Host     string
	User     string
	Password string
	Port     string
	SSLMode  string
}

type NatsConfig struct {
	URL     string
	Subject string
	Cluster string
	Client  string
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config error: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ummarshal to config struct is failed: %w", err)
	}

	return &cfg, nil
}
