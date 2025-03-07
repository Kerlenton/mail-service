package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("test_config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("Expected server port 9090, got %s", cfg.Server.Port)
	}
	if cfg.Database.Host != "testhost" {
		t.Errorf("Expected database host testhost, got %s", cfg.Database.Host)
	}
	if cfg.Redis.Host != "testredis" {
		t.Errorf("Expected redis host testredis, got %s", cfg.Redis.Host)
	}
	if cfg.RabbitMQ.URL != "amqp://test:test@localhost:5672/" {
		t.Errorf("Expected RabbitMQ URL amqp://test:test@localhost:5672/, got %s", cfg.RabbitMQ.URL)
	}
}
