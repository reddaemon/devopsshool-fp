package postgres

import (
	"final-project/internal/config"

	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPsqlDb(c *config.Config) (*pgxpool.Pool, error) {
	fmt.Println(c.Postgres.PostgresqlHost)
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(c.Postgres.PostgresqlUser),
		url.QueryEscape(c.Postgres.PostgresqlPassword),
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlDbname,
		5)

	ctx, _ := context.WithCancel(context.Background())

	//Сконфигурируем пул, задав для него максимальное количество соединений
	poolConfig, _ := pgxpool.ParseConfig(connStr)
	poolConfig.MaxConns = 5

	//Получаем пул соединений, используя контекст и конфиг
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
		os.Exit(1)
	}

	return pool, nil
}
