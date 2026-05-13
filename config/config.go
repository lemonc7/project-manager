package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Logger   LoggerConfig   `yaml:"logger"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port         int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User         string `yaml:"user" env:"DB_USER"`
	Password     string `yaml:"password" env:"DB_PASSWORD"`
	DBname       string `yaml:"dbname" env:"DB_NAME"`
	SSLMode      string `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"disable"`
	MaxIdleConns int    `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"10"`
	MaxOpenConns int    `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"20"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBname, d.SSLMode)
}

type ServerConfig struct {
	Port         int           `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
}

type JWTConfig struct {
	Secret   string        `yaml:"secret" env:"JWT_SECRET"`
	Duration time.Duration `yaml:"duration" env:"JWT_DURATION" env-default:"24h"`
}

type LoggerConfig struct {
	TZ string `yaml:"tz" env:"TZ" env-default:"Asia/Shanghai"`
}

func NewConfig(filePath string) (*Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(filePath, &config); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return &config, nil
}
