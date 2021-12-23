package config

import (
	"errors"
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
	PostgresqlSSLMode  bool
}

// Redis config
type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisUsername string
}

// Load config file from given path
func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath(".")
	if os.Getenv("ENV") == "PRODUCTION" {
		v.SetConfigName(".config")
	} else if os.Getenv("ENV") == "STAGING" {
		v.SetConfigName(".stageconfig")
	} else if os.Getenv("ENV") == "DEV" || os.Getenv("ENV") == "" {
		v.SetConfigName(".devconfig")
	}
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
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
