package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host string
	Port int
}

type RabbitMQConfig struct {
	URL string
}
