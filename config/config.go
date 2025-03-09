package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RabbitMQConfig struct {
	URL string `mapstructure:"url"`
}

type AuthConfig struct {
	JwtSecret string `mapstructure:"jwt_secret"`
}

func LoadConfig(path string) (*Config, error) {
	// Create a new Viper instance for isolation.
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// Enable overriding via environment variables.
	v.SetEnvPrefix("MAIL_SERVICE")
	v.AutomaticEnv()

	// Read configuration file.
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config into struct.
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
