package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   Logger
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  string
}

// Redis config
type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisUsername string
}

// Load config file from given path
func LoadConfig() (*Config, error) {
	v := viper.New()
	var c Config

	v.AddConfigPath(".")
	if os.Getenv("ENV") == "PRODUCTION" {
		v.SetConfigName(".config")
	} else if os.Getenv("ENV") == "STAGING" {
		v.SetConfigName(".stageconfig")
	} else if os.Getenv("ENV") == "DEV" || os.Getenv("ENV") == "" {
		v.SetConfigName(".devconfig")
	}
	v.AutomaticEnv()
	c.Postgres.PostgresqlDbname = v.GetString("POSTGRES_DB_NAME")
	c.Postgres.PostgresqlHost = v.GetString("POSTGRES_HOST")
	c.Postgres.PostgresqlUser = v.GetString("POSTGRES_USER")
	c.Postgres.PostgresqlPassword = v.GetString("POSTGRES_PASSWORD")
	c.Postgres.PostgresqlPort = v.GetString("POSTGRES_PORT")
	c.Postgres.PostgresqlSSLMode = v.GetString("POSTGRES_SSL_MODE")
	c.Redis.RedisHost = v.GetString("REDIS_HOST")
	c.Redis.RedisPort = v.GetString("REDIS_PORT")
	c.Redis.RedisUsername = v.GetString("REDIS_USER")
	c.Redis.RedisPassword = v.GetString("REDIS_PASSWORD")

	return &c, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
