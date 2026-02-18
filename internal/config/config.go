package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Env int

const (
	Local Env = iota
	Development
	Production
)

var envMap = map[string]Env{
	"LOCAL":       Local,
	"DEVELOPMENT": Development,
	"PRODUCTION":  Production,
}

type DB struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" required:"true"`
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	DBName   string `envconfig:"DB_NAME" required:"true"`
}

type Auth struct {
	SecretKey                 []byte `envconfig:"AUTH_SECRET_KEY" required:"true"`
	AccessTokenExpiresMinutes time.Duration
	RefreshTokenExpiresDays   time.Duration
}

type Server struct {
	Port string `envconfig:"SERVER_PORT" required:"true"`
}

type Config struct {
	DB     DB
	Server Server
	Auth   Auth
}

func (cfg *DB) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
}

func Load() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return nil, fmt.Errorf("ENVIRONMENT is not set")
	}

	if envMap[env] == Local {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("loading .env: %w", err)
		}
	}

	var dbCfg DB
	if err := envconfig.Process("", &dbCfg); err != nil {
		return nil, fmt.Errorf("processing db config: %w", err)
	}

	var serverCfg Server
	if err := envconfig.Process("", &serverCfg); err != nil {
		return nil, fmt.Errorf("processing server config: %w", err)
	}

	var authCfg Auth
	if err := envconfig.Process("", &authCfg); err != nil {
		return nil, fmt.Errorf("processing auth config: %w", err)
	}

	if authCfg.AccessTokenExpiresMinutes == 0 {
		authCfg.AccessTokenExpiresMinutes = 60 * time.Minute
	}
	if authCfg.RefreshTokenExpiresDays == 0 {
		authCfg.RefreshTokenExpiresDays = 7 * 24 * time.Hour
	}

	return &Config{
		DB:     dbCfg,
		Server: serverCfg,
		Auth:   authCfg,
	}, nil
}
