package database

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.TODO()
)

type RedisConfig struct {
	Addr     string
	Password string
}

func NewRedisDB(cfg RedisConfig) (*Database, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Database{
		DBType: REDIS,
		Redis:  rdb,
	}, nil
}
