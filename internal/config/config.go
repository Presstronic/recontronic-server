package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Worker   WorkerConfig
	Logging  LoggingConfig
	Security SecurityConfig
}

type ServerConfig struct {
	RESTPort     int
	GRPCPort     int
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type WorkerConfig struct {
	Concurrency int
	QueueName   string
}

type LoggingConfig struct {
	Level  string
	Format string
}

type SecurityConfig struct {
	APIKey         string
	JWTSecret      string
	RateLimitRPS   int
	AllowedOrigins []string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.restport", 8080)
	viper.SetDefault("server.grpcport", 9090)
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.readtimeout", "15s")
	viper.SetDefault("server.writetimeout", "15s")
	viper.SetDefault("server.idletimeout", "60s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.dbname", "recontronic")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.maxopenconns", 25)
	viper.SetDefault("database.maxidleconns", 5)
	viper.SetDefault("database.connmaxlifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// Worker defaults
	viper.SetDefault("worker.concurrency", 10)
	viper.SetDefault("worker.queuename", "recontronic")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")

	// Security defaults
	viper.SetDefault("security.ratelimitrps", 100)
	viper.SetDefault("security.allowedorigins", []string{"*"})
}
