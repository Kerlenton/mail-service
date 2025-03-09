package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	viper.SetConfigType("yaml")
	viper.Set("server.port", "9090")
	viper.Set("database.host", "testhost")
	viper.Set("redis.host", "testredis")
	viper.Set("rabbitmq.url", "amqp://test:test@localhost:5672/")

	cfg, err := LoadConfig("test_config.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "testhost", cfg.Database.Host)
	assert.Equal(t, "testredis", cfg.Redis.Host)
	assert.Equal(t, "amqp://test:test@localhost:5672/", cfg.RabbitMQ.URL)
}
