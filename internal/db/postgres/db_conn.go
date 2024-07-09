package postgres

import (
	"final-project/internal/config"

	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPsqlDb(c *config.Config) (*pgxpool.Pool, error) {
	if os.Getenv("ENV") == "PRODUCTION" {
		connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=%d",
			"postgres",
			url.QueryEscape(c.Postgres.PostgresqlUser),
			url.QueryEscape(c.Postgres.PostgresqlPassword),
			c.Postgres.PostgresqlHost,
			c.Postgres.PostgresqlPort,
			c.Postgres.PostgresqlDbname,
			c.Postgres.PostgresqlSSLMode,
			5)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		//Сконфигурируем пул, задав для него максимальное количество соединений
		poolConfig, _ := pgxpool.ParseConfig(connStr)
		poolConfig.MaxConns = 5

		//Получаем пул соединений, используя контекст и конфиг
		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
			os.Exit(1)
		}

		return pool, nil
	} else {
		connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
			"postgres",
			url.QueryEscape(c.Postgres.PostgresqlUser),
			url.QueryEscape(c.Postgres.PostgresqlPassword),
			c.Postgres.PostgresqlHost,
			c.Postgres.PostgresqlPort,
			c.Postgres.PostgresqlDbname,
			5)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		//Сконфигурируем пул, задав для него максимальное количество соединений
		poolConfig, _ := pgxpool.ParseConfig(connStr)
		poolConfig.MaxConns = 5

		//Получаем пул соединений, используя контекст и конфиг
		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
			os.Exit(1)
		}

		return pool, nil
	}

}
