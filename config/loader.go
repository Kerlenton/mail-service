package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
		return nil, err
	}

	return &config, nil
}
