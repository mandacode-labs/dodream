// Package config provides centralized configuration loading and validation.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	NATS     NATSConfig     `mapstructure:"nats"`
	Engine   EngineConfig   `mapstructure:"engine"`
}

// AppConfig holds application metadata.
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Addr string `mapstructure:"addr"`
}

// DatabaseConfig holds database configuration.
type DatabaseConfig struct {
	URL      string `mapstructure:"url"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	URL string `mapstructure:"url"`
}

// NATSConfig holds NATS configuration.
type NATSConfig struct {
	URL string `mapstructure:"url"`
}

// EngineConfig holds engine-specific configuration.
type EngineConfig struct {
	WindowSize int    `mapstructure:"window_size"`
	GRPCAddr   string `mapstructure:"grpc_addr"`
}

// Load loads configuration with the following precedence:
// 1. Defaults
// 2. Config file (if exists)
// 3. Environment variables
func Load() *Config {
	v := viper.New()
	setDefaults(v)
	bindEnvVars(v)

	v.SetConfigName("dodream")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.dodream")
	v.AddConfigPath("/etc/dodream/")

	_ = v.ReadInConfig() // optional, ignore error if no config file

	var cfg Config
	_ = v.Unmarshal(&cfg)

	return &cfg
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "dodream")
	v.SetDefault("app.version", "dev")
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("database.url", "postgres://dodream:dodream@localhost:5432/dodream?sslmode=disable")
	v.SetDefault("redis.url", "redis://localhost:6379")
	v.SetDefault("nats.url", "nats://localhost:4222")
	v.SetDefault("engine.window_size", 50)
	v.SetDefault("engine.grpc_addr", ":50051")
}

func bindEnvVars(v *viper.Viper) {
	_ = v.BindEnv("server.addr", "SERVER_ADDR")
	_ = v.BindEnv("database.url", "DATABASE_URL")
	_ = v.BindEnv("redis.url", "REDIS_URL")
	_ = v.BindEnv("nats.url", "NATS_URL")
	_ = v.BindEnv("engine.window_size", "ENGINE_WINDOW_SIZE")
	_ = v.BindEnv("engine.grpc_addr", "ENGINE_GRPC_ADDR")
}

// DSN returns the database connection string.
func (c *DatabaseConfig) DSN() string {
	if c.URL != "" {
		return c.URL
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

// RedactedDSN returns the DSN with the password masked for logging.
func (c *DatabaseConfig) RedactedDSN() string {
	dsn := c.DSN()
	// Simple masking: replace password segment
	if idx := strings.Index(dsn, ":"); idx != -1 {
		if atIdx := strings.Index(dsn, "@"); atIdx != -1 {
			return dsn[:idx+1] + "***" + dsn[atIdx:]
		}
	}
	return dsn
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.NATS.URL == "" {
		return fmt.Errorf("nats URL is required")
	}
	if c.Redis.URL == "" {
		return fmt.Errorf("redis URL is required")
	}
	if c.Engine.WindowSize < 1 {
		return fmt.Errorf("engine window size must be at least 1")
	}
	return nil
}
