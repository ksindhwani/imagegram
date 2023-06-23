package config

import (
	"time"

	"github.com/caarlos0/env/v8"
)

const (
	defaultLogLevel             = "info"
	defaultLogFormat            = "json"
	defaultAddr                 = "0.0.0.0:8000"
	defaultRemoteTimeout        = 30 * time.Second
	defaultServerReadTimeout    = 2 * time.Minute
	defaultServerWriteTimeout   = 2 * time.Minute
	defaultServerIdleTimeout    = 5 * time.Minute
	defaultDBMaxIdleConnections = 10
	defaultDBMaxOpenConnections = 5
	defaultDBMaxConnLifeTime    = 30 * time.Minute
	defaultDBUserID             = "root"
	defaultDBPassword           = ""
	defaultDBHostName           = "localhost"
	defaultDBPort               = 3306
	defaultDBDatabaseName       = "database"
)

type Config struct {
	Addr                 string        `env:"ADDR"` // e.g. 0.0.0.0:8000
	LogLevel             string        `env:"LOG_LEVEL" validate:"oneof=debug info warn error fatal panic"`
	LogFormat            string        `env:"LOG_FORMAT" validate:"oneof=text console json"`
	ServerReadTimeout    time.Duration `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout   time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout    time.Duration `env:"SERVER_IDLE_TIMEOUT"`
	DBUserID             string        `env:"DB_USER"`
	DBPassword           string        `env:"DB_PASSWORD"`
	DBHostName           string        `env:"DB_HOST"`
	DBPort               int           `env:"DB_PORT"`
	DBDatabaseName       string        `env:"DB_NAME"`
	DBMaxIdleConnections int           `env:"DB_MAX_IDLE_CONNECTIONS"`
	DBMaxOpenConnections int           `env:"DB_MAX_OPEN_CONNECTIONS"`
	DBMaxConnLifetime    time.Duration `env:"DB_MAX_CONN_LIFETIME"`
}

func New() (*Config, error) {
	cfg := Config{
		Addr:                 defaultAddr,
		LogLevel:             defaultLogLevel,
		LogFormat:            defaultLogFormat,
		ServerReadTimeout:    defaultServerReadTimeout,
		ServerWriteTimeout:   defaultServerWriteTimeout,
		ServerIdleTimeout:    defaultServerIdleTimeout,
		DBUserID:             defaultDBUserID,
		DBPassword:           defaultDBPassword,
		DBHostName:           defaultDBHostName,
		DBPort:               defaultDBPort,
		DBDatabaseName:       defaultDBDatabaseName,
		DBMaxIdleConnections: defaultDBMaxIdleConnections,
		DBMaxOpenConnections: defaultDBMaxOpenConnections,
		DBMaxConnLifetime:    defaultDBMaxConnLifeTime,
	}
	// load .env file
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
