package redis

import (
	"context"
	"final-project/internal/config"

	"github.com/go-redis/redis/v8"
)

func NewRedisConn(c *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis.RedisHost + ":" + c.Redis.RedisPort,
		//Username: c.Redis.RedisUsername,
		Password: c.Redis.RedisPassword,
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
